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
	return ReadHelper.Require(klabels.VariableKind, analysis.Labels, func(label klabels.Label) (string, error) {
		variables := (label).Data.([]klabels.VariableDeclaration)
		return ReadHelper.GetParameters("+k:env", analysis.Labels, func(parameters []string) (string, error) {
			var kind string
			for i, variable := range variables {
				if variable.Value != nil && kstrings.HasPrefix(*variable.Value, "\"{$ENV:") && kstrings.HasSuffix(*variable.Value, "}\"") {
					environmentKey := strings.TrimPrefix(*variable.Value, "\"{$ENV:")
					environmentKey = strings.TrimSuffix(environmentKey, "}\"")
					environmentValue := os.Getenv(environmentKey)

					if len(parameters) == 1 {
						kind = parameters[0]
						variable.Type = &parameters[0]
					} else {
						if variable.Type != nil {
							kind = *variable.Type
						}
					}

					newValue := SyntaxHelper.AutoValueBasedOnKind(environmentValue, kind)
					variable.Value = &newValue
					variables[i] = variable
				}
			}
			return SyntaxHelper.AutoVariableDeclaration(line, index, stack, variables), nil
		})
	})
}
