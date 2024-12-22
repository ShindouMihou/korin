package korin

import (
	"slices"
	"strconv"
	"strings"
)

type Generator int

var WriteAssistant Generator = 0

func (generator Generator) Func(name string, parameters [][]string, results []string) string {
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

func (generator Generator) If(statement string) string {
	return "if " + statement
}

func (generator Generator) OpenBracket() string {
	return "{"
}

func (generator Generator) CloseBracket() string {
	return "}"
}

func (generator Generator) Else(statement *string) string {
	if statement != nil {
		return "else " + *statement
	}

	return "else"
}

func (generator Generator) NewLine() string {
	return "\n"
}

func (generator Generator) VariableDeclaration(reassign bool, names []string, values []string) string {
	text := generator.CommaSeparate(names)

	if reassign {
		text += " = "
	} else {
		text += " := "
	}

	text += generator.CommaSeparate(values)
	return text
}

func (generator Generator) CommaSeparate(parts []string) string {
	text := ""
	for i, part := range parts {
		if i > 0 {
			text += ", "
		}
		text += part
	}
	return text
}

func (generator Generator) Return(values []string) string {
	return "return " + generator.CommaSeparate(values)
}

func (generator Generator) TabSize(len int) string {
	return strings.Repeat("\t", len)
}

func (generator Generator) Tab() string {
	return "\t"
}

func (generator Generator) Quote(text string) string {
	return strconv.Quote(text)
}

func (generator Generator) TabSizeFrom(line string) string {
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

func (generator Generator) Call(pkg string, method string, parameters []string) string {
	if pkg == "" {
		return method + "(" + generator.CommaSeparate(parameters) + ")"
	}
	return pkg + "." + method + "(" + generator.CommaSeparate(parameters) + ")"
}

func (generator Generator) Import(packages []string) string {
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
