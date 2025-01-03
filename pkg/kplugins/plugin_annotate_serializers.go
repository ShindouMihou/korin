package kplugins

import (
	"errors"
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/korin/pkg/klabels"
)

type PluginSerializerAnnotations struct {
	Plugin
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

func (p PluginSerializerAnnotations) Process(line string, index int, headers *Headers, stack []klabels.Analysis) (string, error) {
	analysis := stack[index]
	fieldDeclaration := ReadHelper.Get(klabels.FieldDeclarationKind, analysis.Labels)

	if fieldDeclaration != nil {
		field := (*fieldDeclaration).Data.(klabels.FieldDeclaration)

		shouldAnnotate, parameters := ReadHelper.Parameters("+k:named", analysis.Labels)
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

	return "", nil
}
