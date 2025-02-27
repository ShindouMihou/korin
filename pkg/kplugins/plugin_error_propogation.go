package kplugins

import (
	"errors"
	"fmt"
	"github.com/ShindouMihou/korin/pkg/klabels"
	"strconv"
	"strings"
)

type ErrorPropogationPlugin struct {
	Plugin
}

func NewErrorPropogationPlugin() ErrorPropogationPlugin {
	return ErrorPropogationPlugin{}
}

func (p ErrorPropogationPlugin) Name() string {
	return "KorinErrorPropogation"
}

func (p ErrorPropogationPlugin) Group() string {
	return "pw.mihou"
}

func (p ErrorPropogationPlugin) Version() string {
	return "1.0.0"
}

func (p ErrorPropogationPlugin) Context(file string) *any {
	return nil
}

func (p ErrorPropogationPlugin) FreeContext(file string) {
}

func (p ErrorPropogationPlugin) Process(line string, index int, headers *Headers, stack []klabels.Analysis, context *any) (string, error) {
	analysis := stack[index]

	varDeclaration := ReadHelper.Get(klabels.VariableKind, analysis.Labels)
	shouldPropogate, parameters := ReadHelper.Parameters("+k:float", analysis.Labels)

	if shouldPropogate {
		if varDeclaration == nil {
			return NoChanges, fmt.Errorf("cannot float error in a non-variable-declaration line (Line %d)", analysis.Line+1)
		}

		variables := (*varDeclaration).Data.([]klabels.VariableDeclaration)
		var errorVar *klabels.VariableDeclaration

		if len(parameters) > 0 {
			pos, err := strconv.Atoi(parameters[0])
			if err != nil {
				return NoChanges, errors.Join(
					fmt.Errorf("expected `int` from k:float's first parameter (position of error variable) "+
						"(Line %d)", analysis.Line+1),
					err,
				)
			}
			if len(variables) < pos {
				return NoChanges, fmt.Errorf("no variable in the position of %d in k:float (Line %d)", pos, analysis.Line+1)
			}
			errorVar = &variables[pos]
		} else {
			for _, variable := range variables {
				variable := variable
				if variable.Name == "_" || variable.Name == "err" {
					errorVar = &variable
					break
				}
			}
		}

		if errorVar == nil {
			return NoChanges, fmt.Errorf("no error variable provided in k:float (Line %d): "+
				"use _ or `err` for the error name, or provide the position as an argument", analysis.Line+1)
		}

		var function *klabels.Label
		for index := analysis.Line; function == nil && index > 0; index-- {
			fn := ReadHelper.Get(klabels.FunctionKind, stack[index].Labels)
			if fn != nil {
				function = fn
				break
			}
		}

		if function == nil {
			return NoChanges, fmt.Errorf("k:float variable is declared in a non-function scope. (Line %d)", analysis.Line+1)
		}

		results := function.Data.(klabels.FunctionDeclaration).Result
		var resultValues []string
		for _, result := range results {
			result := result
			if result == "string" {
				resultValues = append(resultValues, "\"\"")
			} else if result == "bool" {
				resultValues = append(resultValues, "false")
			} else if result == "rune" {
				resultValues = append(resultValues, "''")
			} else if strings.Contains(result, "float") || result == "byte" || strings.Contains(result, "int") {
				resultValues = append(resultValues, "0")
			} else if result == "error" {
				name := (*errorVar).Name
				if name == "_" {
					name = "err"
				}
				resultValues = append(resultValues, name)
			} else {
				resultValues = append(resultValues, "nil")
			}
		}

		var declarationNames []string
		var declarationValues []string

		var errorVariableName string
		for _, variable := range variables {
			variable := variable
			if variable == *errorVar {
				if variable.Name == "_" {
					errorVariableName = "err"
					declarationNames = append(declarationNames, "err")
				} else {
					errorVariableName = variable.Name
				}
			} else {
				declarationNames = append(declarationNames, variable.Name)
			}

			if variable.Value != nil {
				declarationValues = append(declarationValues, *variable.Value)
			}
		}

		tab := SyntaxHelper.TabSizeFrom(line)

		line := tab
		line += SyntaxHelper.VariableDeclaration(errorVar.Reassignment, declarationNames, declarationValues, "")
		line += SyntaxHelper.NewLine()
		line += tab + SyntaxHelper.If(errorVariableName+" != nil") + " " + SyntaxHelper.OpenBracket()
		line += SyntaxHelper.NewLine()
		line += tab + SyntaxHelper.Tab() + SyntaxHelper.Return(resultValues)
		line += SyntaxHelper.NewLine()
		line += tab + SyntaxHelper.CloseBracket()
		return line, nil
	}

	return NoChanges, nil
}
