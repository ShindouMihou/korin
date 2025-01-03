package kplugins

import (
	"github.com/ShindouMihou/korin/pkg/klabels"
	"slices"
	"strconv"
	"strings"
)

type SyntaxHelperType int

var SyntaxHelper SyntaxHelperType = 0

// Func generates a function with the given name, parameters and results, this doesn't include the opening and closing bracket.
func (generator SyntaxHelperType) Func(name string, parameters [][]string, results []string) string {
	text := "func " + name + "("
	if len(parameters) > 0 {
		for index, parameter := range parameters {
			if len(parameter) < 2 {
				panic("a parameter must be two strings, [name, type]")
			}

			if index != 0 {
				text += ", "
			}

			text += parameter[0]
			text += " "
			text += parameter[1]
		}
	}
	text += ")"

	if len(results) > 0 {
		if len(results) == 1 {
			text += " " + results[0] + ""
		} else {
			text += " ("
			for index, t := range results {
				if index != 0 {
					text += ", "
				}

				text += t
			}
			text += ")"
		}
	}

	return text
}

// If generates an if statement with the given statement, this doesn't include the opening and closing bracket.
func (generator SyntaxHelperType) If(statement string) string {
	return "if " + statement
}

// OpenBracket returns the opening bracket.
func (generator SyntaxHelperType) OpenBracket() string {
	return "{"
}

// CloseBracket returns the closing bracket.
func (generator SyntaxHelperType) CloseBracket() string {
	return "}"
}

// Else generates an else statement with the given statement, this doesn't include the opening and closing bracket.
func (generator SyntaxHelperType) Else(statement *string) string {
	if statement != nil {
		return "else " + *statement
	}

	return "else"
}

// NewLine returns a new line.
func (generator SyntaxHelperType) NewLine() string {
	return "\n"
}

// VariablesDeclaration generates a variable declaration given the variables, whether it's a reassignment, and its kind.
// Unlike VariableDeclaration, this will automatically generate the names and values based on the variables.
func (generator SyntaxHelperType) VariablesDeclaration(reassign bool, variables []klabels.VariableDeclaration, kind string) string {
	var declarationNames []string
	var declarationValues []string

	for _, variable := range variables {
		variable := variable
		declarationNames = append(declarationNames, strings.TrimSpace(variable.Name))
		if variable.Value != nil {
			declarationValues = append(declarationValues, *variable.Value)
		}
	}

	return generator.VariableDeclaration(reassign, declarationNames, declarationValues, kind)
}

// VariableDeclaration generates a variable declaration given the names, values, whether it's a reassignment, and its kind.
// Unlike VariablesDeclaration, this is more manual and will allow you to go finer with the details.
func (generator SyntaxHelperType) VariableDeclaration(reassign bool, names []string, values []string, kind string) string {
	text := generator.CommaSeparate(names)

	if kind != "" {
		text += " " + kind
	}

	if reassign {
		text += " = "
	} else {
		text += " := "
	}

	text += generator.CommaSeparate(values)
	return text
}

// CommaSeparate will comma separate the given parts.
func (generator SyntaxHelperType) CommaSeparate(parts []string) string {
	text := ""
	for i, part := range parts {
		if i > 0 {
			text += ", "
		}
		text += part
	}
	return text
}

// Return generates a return statement with the given values.
func (generator SyntaxHelperType) Return(values []string) string {
	return "return " + generator.CommaSeparate(values)
}

// TabSize returns a tab size based on the given length.
func (generator SyntaxHelperType) TabSize(len int) string {
	return strings.Repeat("\t", len)
}

// Tab returns a tab character.
func (generator SyntaxHelperType) Tab() string {
	return "\t"
}

// Quote returns a quoted string, this will ensure that it is properly escaped.
func (generator SyntaxHelperType) Quote(text string) string {
	return strconv.Quote(text)
}

// SingleQuote returns a single quoted string, this will ensure that it is properly escaped.
func (generator SyntaxHelperType) SingleQuote(text string) string {
	return strconv.QuoteRune(rune(text[0]))
}

// TabSizeFrom returns a tab size based on the given line's tab spacing.
func (generator SyntaxHelperType) TabSizeFrom(line string) string {
	tabCount := 0
	for _, char := range line {
		if char == '\t' {
			tabCount++
		} else {
			break
		}
	}

	return generator.TabSize(tabCount)
}

// Call generates a function call with the given package, method and parameters.
func (generator SyntaxHelperType) Call(pkg string, method string, parameters []string) string {
	if pkg == "" {
		return method + "(" + generator.CommaSeparate(parameters) + ")"
	}
	return pkg + "." + method + "(" + generator.CommaSeparate(parameters) + ")"
}

// Import generates an import statement with the given packages, this will auto-handle
// multi-line imports and single-line imports.
func (generator SyntaxHelperType) Import(packages []string) string {
	if len(packages) == 1 {
		return "import " + strconv.Quote(packages[0])
	} else {
		text := "import ("
		text += generator.NewLine()

		slices.Sort(packages)

		for index, pkg := range packages {
			if index > 0 {
				text += generator.NewLine()
			}
			text += generator.Tab() + strconv.Quote(pkg)
		}
		text += generator.NewLine()
		text += ")"
		return text
	}
}

// TypeDeclaration generates a type declaration with the given name and kind, this doesn't include the
// opening and closing brackets.
func (generator SyntaxHelperType) TypeDeclaration(name string, kind string) string {
	return "type " + name + " " + kind
}

// FieldDeclaration generates a field declaration with the given name, kind and annotation.
func (generator SyntaxHelperType) FieldDeclaration(name string, kind string, annotation string) string {
	text := name + " " + kind
	if annotation != "" {
		text += " `" + annotation + "`"
	}
	return text
}

// AutoValueBasedOnKind returns the value based on the kind of the variable, it will automatically
// assume String if the type is unknown.
func (generator SyntaxHelperType) AutoValueBasedOnKind(value string, kind string) string {
	if kind == "rune" {
		return SyntaxHelper.SingleQuote(value)
	} else if strings.Contains(kind, "float") || kind == "byte" || strings.Contains(kind, "int") || kind == "bool" {
		return value
	} else {
		return SyntaxHelper.Quote(value)
	}
}

// AutoVariableDeclaration will automatically generate the variable declaration based on the kind of the variable, this will
// handle cases such as const, var, multi-line const or var and reassignment variable declarations.
func (generator SyntaxHelperType) AutoVariableDeclaration(line string, index int, stack []klabels.Analysis, variables []klabels.VariableDeclaration) string {
	labels := stack[index].Labels

	var constDeclaration, varDeclaration *klabels.Label
	ReadHelper.Filter(labels, func(label klabels.Label) {
		switch label.Kind {
		case klabels.ConstDeclarationKind:
			constDeclaration = &label
		case klabels.VarDeclarationKind:
			varDeclaration = &label
		}
	}, klabels.ConstDeclarationKind, klabels.VarDeclarationKind)
	isMultiLineConstOrVar := AnalysisHelper.CheckMultiLineConstOrVar(index, stack)

	isReassignment := variables[0].Reassignment

	var kind string
	if vt := variables[0].Type; vt != nil {
		kind = *vt
	}

	var prefix string
	if (constDeclaration != nil || varDeclaration != nil) && !isMultiLineConstOrVar {
		isReassignment = true
		if constDeclaration != nil {
			prefix = "const "
		} else {
			prefix = "var "
		}
	}

	if isMultiLineConstOrVar {
		isReassignment = true
	}

	newLine := SyntaxHelper.TabSizeFrom(line)
	newLine += prefix + SyntaxHelper.VariablesDeclaration(isReassignment, variables, kind)
	return newLine
}
