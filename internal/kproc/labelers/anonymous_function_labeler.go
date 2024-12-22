package labelers

import (
	"korin/pkg/klabels"
	"strings"
)

func AnonymousFunction(line string) klabels.Label {
	label := klabels.Label{Kind: klabels.FunctionKind}

	declaration := line[strings.Index(line, "func"):]
	var parameters []klabels.FunctionParameter
	var resultTypes []string

	currentToken := ""
	parameterName := ""
	currentLocation := "parameters"

	for _, char := range declaration[4:] {
		if char == ',' {
			if currentLocation == "parameters" && parameterName != "" && currentToken != "" {
				parameters = append(parameters, klabels.FunctionParameter{Name: parameterName, Type: currentToken})
				parameterName = ""

				currentToken = ""
				continue
			}
			if currentLocation == "result" {
				resultTypes = append(resultTypes, currentToken)
				currentToken = ""
				continue
			}
		}
		if char == ' ' {
			if currentLocation == "result" && currentToken != "" {
				resultTypes = append(resultTypes, currentToken)
				currentToken = ""

				currentLocation = "outside"
				continue
			}
			if currentLocation == "parameters" && parameterName == "" && currentToken != "" {
				parameterName = currentToken
				currentToken = ""
				continue
			}
			continue
		}
		if char == '(' {
			continue
		}

		if char == ')' {
			if currentLocation == "parameters" {
				if currentToken != "" {
					parameters = append(parameters, klabels.FunctionParameter{Name: parameterName, Type: currentToken})
					parameterName = ""

					currentToken = ""
				}

				currentLocation = "result"
				continue
			}

			if currentLocation == "result" {
				resultTypes = append(resultTypes, currentToken)
				currentToken = ""

				currentLocation = "outside"
				continue
			}
		}
		currentToken = currentToken + string(char)
	}

	label.Data = klabels.FunctionDeclaration{Name: "", Parameters: parameters, Result: resultTypes}
	return label
}
