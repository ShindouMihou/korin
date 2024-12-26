package labelers

import "github.com/ShindouMihou/korin/pkg/klabels"

func FuncLine(line string) klabels.Label {
	label := klabels.Label{Kind: klabels.FunctionKind}

	currentName := ""

	var parameters []klabels.FunctionParameter
	var resultTypes []string

	currentToken := ""
	parameterName := ""
	currentLocation := "name"

	for _, char := range line[5:] {
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
			if currentLocation == "name" {
				currentName = currentToken
				currentToken = ""

				currentLocation = "parameters"
				continue
			}
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

	label.Data = klabels.FunctionDeclaration{Name: currentName, Parameters: parameters, Result: resultTypes}
	return label
}
