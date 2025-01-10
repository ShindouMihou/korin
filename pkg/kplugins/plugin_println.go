package kplugins

import (
	"errors"
	"fmt"
	"github.com/ShindouMihou/go-little-utils/slices"
	"github.com/ShindouMihou/korin/pkg/klabels"
	"strconv"
	"strings"
)

type PrintLinePlugin struct {
	Plugin
}

func NewPrintLinePlugin() PrintLinePlugin {
	return PrintLinePlugin{}
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

func (p PrintLinePlugin) Context(file string) *any {
	return nil
}

func (p PrintLinePlugin) FreeContext(file string) {
}

func (p PrintLinePlugin) Process(line string, index int, headers *Headers, stack []klabels.Analysis, context *any) (string, error) {
	analysis := stack[index]

	varDeclaration := ReadHelper.Get(klabels.VariableKind, analysis.Labels)
	shouldPrint, parameters := ReadHelper.Parameters("+k:println", analysis.Labels)

	if shouldPrint {
		if varDeclaration == nil {
			return NoChanges, fmt.Errorf("cannot print non-variable declaration (Line %d)", analysis.Line+1)
		}

		variables := (*varDeclaration).Data.([]klabels.VariableDeclaration)
		var printVariables []klabels.VariableDeclaration

		if len(parameters) > 0 {
			for index, parameter := range parameters {
				pos, err := strconv.Atoi(parameter)
				if err != nil {
					return NoChanges, errors.Join(
						fmt.Errorf("expected `int` from k:println param %d (position of print variable) "+
							"(Line %d)", index, analysis.Line+1),
						err,
					)
				}
				if len(variables) < pos {
					return NoChanges, fmt.Errorf("no variable in the position of %d in k:println (Line %d)", pos, analysis.Line+1)
				}
				printVariables = append(printVariables, variables[pos])
			}
		} else {
			printVariables = slices.AllMatches(variables, func(b klabels.VariableDeclaration) bool {
				return b.Name != "_"
			})
		}

		tab := SyntaxHelper.TabSizeFrom(line)

		headers.Import("fmt")
		comment := ReadHelper.Get(klabels.CommentKind, analysis.Labels)
		line := strings.TrimSuffix(line, "// "+comment.Data.(string))
		for _, variable := range printVariables {
			line += SyntaxHelper.NewLine()
			line += tab + SyntaxHelper.Call("fmt", "Println", []string{
				SyntaxHelper.Quote(variable.Name+": ") + ", " + variable.Name,
			})
		}
		line += ""
		return line, nil
	}

	return NoChanges, nil
}
