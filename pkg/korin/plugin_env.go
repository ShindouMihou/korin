package korin

import (
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/korin/pkg/klabels"
	"os"
	"strings"
)

type PluginEnvironmentKey struct {
	Plugin
}

func (p PluginEnvironmentKey) Name() string {
	return "KorinEnvironmentKey"
}

func (p PluginEnvironmentKey) Group() string {
	return "pw.mihou"
}

func (p PluginEnvironmentKey) Version() string {
	return "1.0.0"
}

func (p PluginEnvironmentKey) Process(line string, index int, headers *Headers, stack []klabels.Analysis) (string, error) {
	analysis := stack[index]
	variableDeclaration := ReadAssistant.Get(klabels.VariableKind, analysis.Labels)

	constDeclaration := ReadAssistant.Get(klabels.ConstDeclarationKind, analysis.Labels)
	varDeclaration := ReadAssistant.Get(klabels.VarDeclarationKind, analysis.Labels)

	if variableDeclaration != nil {
		variables := (*variableDeclaration).Data.([]klabels.VariableDeclaration)
		shouldAnnotate, parameters := ReadAssistant.Parameters("+k:env", analysis.Labels)
		if shouldAnnotate {
			hasScopeBegin, hasScopeEnd := false, false
			for index := analysis.Line; hasScopeBegin == false && index > 0; index-- {
				for _, label := range stack[index].Labels {
					label := label
					if label.Kind == klabels.ConstScopeBeginKind || label.Kind == klabels.VarScopeBeginKind {
						hasScopeBegin = true
					}
					if label.Kind == klabels.ConstScopeEndKind || label.Kind == klabels.VarScopeEndKind {
						hasScopeEnd = true
						if !hasScopeBegin {
							break
						}
					}
					if !hasScopeBegin && hasScopeEnd {
						break
					}
				}
				if !hasScopeBegin && hasScopeEnd {
					break
				}
			}

			isReassignment := false
			tab := WriteAssistant.TabSizeFrom(line)
			kind := ""
			for i, variable := range variables {
				if variable.Value != nil && kstrings.HasPrefix(*variable.Value, "\"{$ENV:") && kstrings.HasSuffix(*variable.Value, "}\"") {
					environmentKey := strings.TrimPrefix(*variable.Value, "\"{$ENV:")
					environmentKey = strings.TrimSuffix(environmentKey, "}\"")
					environmentValue := os.Getenv(environmentKey)

					newValue := `"` + environmentValue + `"`
					if len(parameters) == 1 {
						kind = parameters[0]
					} else {
						if variable.Type != nil {
							kind = *variable.Type
						}
					}

					if kind == "rune" {
						newValue = "'" + environmentValue + "'"
						variable.Value = &newValue
					} else if strings.Contains(kind, "float") || kind == "byte" || strings.Contains(kind, "int") || kind == "bool" {
						variable.Value = &environmentValue
					} else {
						variable.Value = &newValue
					}

					variables[i] = variable
				}
				isReassignment = variable.Reassignment
			}

			prefix := ""
			if (constDeclaration != nil || varDeclaration != nil) && !hasScopeBegin {
				isReassignment = true
				if constDeclaration != nil {
					prefix = "const "
				} else {
					prefix = "var "
				}
			}

			if hasScopeBegin {
				isReassignment = true
			}

			var declarationNames []string
			var declarationValues []string
			for _, variable := range variables {
				variable := variable
				declarationNames = append(declarationNames, strings.TrimSpace(variable.Name))
				if variable.Value != nil {
					declarationValues = append(declarationValues, *variable.Value)
				}
			}

			line := tab
			line += prefix + WriteAssistant.VariableDeclaration(isReassignment, declarationNames, declarationValues, kind)
			return line, nil
		}
		return "", nil
	}
	return "", nil
}
