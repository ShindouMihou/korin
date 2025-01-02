package labelers

import (
	"github.com/ShindouMihou/korin/pkg/klabels"
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

	hasNameBeenAppended := false
	hasPriorComma := false

	kind := ""
	for index, char := range line {
		if (index + 1) == len(line) {
			if currentLocation == "type" {
				kind = strings.TrimSpace(currentToken)
				break
			}

			currentToken += string(char)

			variableValues = append(variableValues, strings.TrimSpace(currentToken))
			currentToken = ""
			break
		}

		if char == ',' {
			if currentLocation == "names" {
				variableNames = append(variableNames, strings.TrimSpace(currentToken))

				currentToken = ""

				hasNameBeenAppended = true
				hasPriorComma = true
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

		if char == '=' && !isInsideQuotation && (currentLocation == "names" || currentLocation == "type") {
			if line[index-1] == ':' {
				currentToken = strings.TrimSuffix(currentToken, ":")
			}
			if currentLocation == "type" {
				kind = strings.TrimSpace(currentToken)
			}

			currentToken = ""
			currentLocation = "values"
			continue
		}

		if char == ' ' && currentLocation == "names" {
			if !hasPriorComma && hasNameBeenAppended {
				variableNames = append(variableNames, strings.TrimSpace(currentToken))

				currentToken = ""
				currentLocation = "type"
			} else if !hasNameBeenAppended {
				variableNames = append(variableNames, strings.TrimSpace(currentToken))
				hasNameBeenAppended = false

				currentToken = ""
				currentLocation = "type"
			} else if hasPriorComma {
				hasPriorComma = false
			}
			continue
		}

		currentToken += string(char)

		if index == len(line)-1 && currentLocation == "type" {
			kind = strings.TrimSpace(currentToken)
		}
	}

	for i, name := range variableNames {
		var value *string
		if len(variableValues) > i {
			value = &variableValues[i]
		}
		declaredVariables = append(declaredVariables, klabels.VariableDeclaration{
			Name: name, Value: value, Type: &kind, Reassignment: reassignment})
	}

	label.Data = declaredVariables
	return label
}
