package kplugins

import (
	"errors"
	ptr "github.com/ShindouMihou/go-little-utils"
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/korin/pkg/klabels"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type serializerAnnotationContext map[string]any
type PluginSerializerAnnotations struct {
	Plugin
	contexts cmap.ConcurrentMap[string, serializerAnnotationContext]
}

func NewPluginSerializerAnnotations() PluginSerializerAnnotations {
	return PluginSerializerAnnotations{
		contexts: cmap.New[serializerAnnotationContext](),
	}
}

func (p PluginSerializerAnnotations) Name() string {
	return "KorinSerializerAnnotations"
}

func (p PluginSerializerAnnotations) Group() string {
	return "pw.mihou"
}

func (p PluginSerializerAnnotations) Version() string {
	return "1.0.0"
}

func (p PluginSerializerAnnotations) Context(file string) *any {
	p.contexts.SetIfAbsent(file, serializerAnnotationContext{})
	if ctx, ok := p.contexts.Get(file); ok {
		return ptr.Ptr(any(ctx))
	}
	return p.Context(file)
}

func (p PluginSerializerAnnotations) FreeContext(file string) {
	p.contexts.Remove(file)
}

func (p PluginSerializerAnnotations) Process(line string, index int, headers *Headers, stack []klabels.Analysis, context *any) (string, error) {
	analysis := stack[index]
	ctx := (*context).(serializerAnnotationContext)

	shouldAnnotate, parameters := ReadHelper.Parameters("+k:named", analysis.Labels)
	if shouldAnnotate {
		if len(parameters) <= 0 {
			return NoChanges, errors.New("no parameters found for +k:named, expected at least one serializer (e.g. json)")
		}
		ctx["parameters"] = parameters
		if len(analysis.Labels) == 1 {
			return SkipLine, nil
		}
	}

	typeDeclaration := ReadHelper.Get(klabels.TypeDeclarationKind, analysis.Labels)
	if typeDeclaration != nil && !AnalysisHelper.HasClosingBracket(analysis.Labels) {
		if _, ok := ctx["parameters"]; ok {
			ctx["shouldAnnotate"] = true

			declaration := typeDeclaration.Data.(klabels.TypeDeclaration)
			hasOpeningBracket := AnalysisHelper.HasOpenBracket(analysis.Labels)
			return SyntaxHelper.TypeDeclaration(declaration.Name, declaration.Kind, hasOpeningBracket, false), nil
		}
		return NoChanges, nil
	}

	if _, ok := ctx["shouldAnnotate"]; ok && AnalysisHelper.HasClosingBracket(analysis.Labels) {
		delete(ctx, "shouldAnnotate")
		return NoChanges, nil
	}

	fieldDeclaration := ReadHelper.Get(klabels.FieldDeclarationKind, analysis.Labels)
	if fieldDeclaration != nil {
		field := (*fieldDeclaration).Data.(klabels.FieldDeclaration)

		shouldAnnotate, parameters := ReadHelper.Parameters("+k:named", analysis.Labels)
		if _, ok := ctx["shouldAnnotate"]; ok && !shouldAnnotate {
			shouldAnnotate = true
			parameters = ctx["parameters"].([]string)
		}

		if shouldAnnotate {
			if len(parameters) <= 0 {
				return "", errors.New("no parameters found for +k:named, expected at least one serializer (e.g. json)")
			}

			casingStyle := "snake_case"
			startingIndex := 0
			if parameters[0] == "camelCase" || parameters[0] == "snake_case" || parameters[0] == "original" {
				if len(parameters) <= 1 {
					return line, errors.New("no serializer found for +k:named, expected at least one serializer (e.g. json)")
				}

				if parameters[0] == "camelCase" {
					casingStyle = "camelCase"
				}

				if parameters[0] == "original" {
					casingStyle = "original"
				}

				startingIndex = 1
			}

			for _, serializer := range parameters[startingIndex:] {
				name := field.Name
				if casingStyle == "camelCase" {
					name = kstrings.ToCamelCase(name)
				} else if casingStyle == "original" {
					name = field.Name
				} else {
					name = kstrings.ToSnakeCase(name)
				}

				if len(field.Annotations) > 0 {
					field.Annotations += " "
				}
				field.Annotations += serializer + `:"` + name + `"`
			}

			line = SyntaxHelper.TabSizeFrom(line)
			line += SyntaxHelper.FieldDeclaration(field.Name, field.Type, field.Annotations)
			return line, nil
		}
	}

	return NoChanges, nil
}
