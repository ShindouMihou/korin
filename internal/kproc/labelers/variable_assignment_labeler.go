package labelers

import (
	"korin/pkg/klabels"
	"strings"
)

func VariableAssignment(reassignment bool, line string) klabels.Label {
	label := klabels.Label{Kind: klabels.VariableKind}

	var declaredVariables []klabels.VariableDeclaration

	var variableNames []string
	var variableValues []string

	currentToken := ""

	currentLocation := "names"

	isInsideQuotation := false
	for index, char := range line {
		if (index + 1) == len(line) {
			currentToken += string(char)

			variableValues = append(variableValues, strings.TrimSpace(currentToken))
			currentToken = ""
			break
		}

		if char == ',' {
			if currentLocation == "names" {
				variableNames = append(variableNames, strings.TrimSpace(currentToken))

				currentToken = ""
				continue
			} else if currentLocation == "values" {
				variableValues = append(variableValues, strings.TrimSpace(currentToken))

				currentToken = ""
				continue
			}

			if isInsideQuotation {
				currentToken += string(char)
			}
			continue
		}

		if char == '"' && line[index-1] != '\\' {
			isInsideQuotation = !isInsideQuotation
			currentToken += string(char)
			continue
		}

		// the next parts will be a comment
		if char == '/' && !isInsideQuotation {
			variableValues = append(variableValues, strings.TrimSpace(currentToken))

			currentToken = ""
			break
		}

		if char == '=' && !isInsideQuotation && currentLocation == "names" {
			if line[index-1] == ':' {
				currentToken = strings.TrimSuffix(currentToken, ":")
			}
			variableNames = append(variableNames, strings.TrimSpace(currentToken))

			currentToken = ""
			currentLocation = "values"
			continue
		}

		if char == ' ' && currentLocation == "names" {
			continue
		}

		currentToken += string(char)
	}

	for i, name := range variableNames {
		var value *string
		if len(variableValues) > i {
			value = &variableValues[i]
		}
		declaredVariables = append(declaredVariables, klabels.VariableDeclaration{Name: name, Value: value, Reassignment: reassignment})
	}

	label.Data = declaredVariables
	return label
}
