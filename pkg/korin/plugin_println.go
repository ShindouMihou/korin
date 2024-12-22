package korin

import (
	"errors"
	"fmt"
	"github.com/ShindouMihou/go-little-utils/slices"
	"korin/pkg/klabels"
	"strconv"
	"strings"
)

type PrintLinePlugin struct {
	Plugin
}

func (p PrintLinePlugin) Name() string {
	return "KorinPrintLinePlugin"
}

func (p PrintLinePlugin) Group() string {
	return "pw.mihou"
}

func (p PrintLinePlugin) Version() string {
	return "1.0.0"
}

func (p PrintLinePlugin) Process(line string, index int, headers *Headers, stack []klabels.Analysis) (string, error) {
	analysis := stack[index]

	varDeclaration := ReadAssistant.Get(klabels.VariableKind, analysis.Labels)
	shouldPrint, parameters := ReadAssistant.Parameters("+k:println", analysis.Labels)

	if shouldPrint {
		if varDeclaration == nil {
			return "", fmt.Errorf("cannot print non-variable declaration (Line %d)", analysis.Line+1)
		}

		variables := (*varDeclaration).Data.([]klabels.VariableDeclaration)
		var printVariables []klabels.VariableDeclaration

		if len(parameters) > 0 {
			for index, parameter := range parameters {
				pos, err := strconv.Atoi(parameter)
				if err != nil {
					return "", errors.Join(
						fmt.Errorf("expected `int` from k:println param %d (position of print variable) "+
							"(Line %d)", index, analysis.Line+1),
						err,
					)
				}
				if len(variables) < pos {
					return "", fmt.Errorf("no variable in the position of %d in k:println (Line %d)", pos, analysis.Line+1)
				}
				printVariables = append(printVariables, variables[pos])
			}
		} else {
			printVariables = slices.AllMatches(variables, func(b klabels.VariableDeclaration) bool {
				return b.Name != "_"
			})
		}

		tab := WriteAssistant.TabSizeFrom(line)

		headers.Import("fmt")
		comment := ReadAssistant.Get(klabels.CommentKind, analysis.Labels)
		line := strings.TrimSuffix(line, "// "+comment.Data.(string))
		for _, variable := range printVariables {
			line += WriteAssistant.NewLine()
			line += tab + WriteAssistant.Call("fmt", "Println", []string{
				WriteAssistant.Quote(variable.Name+": ") + ", " + variable.Name,
			})
		}
		line += ""
		return line, nil
	}

	return "", nil
}
