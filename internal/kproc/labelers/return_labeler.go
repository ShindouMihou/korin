package labelers

import (
	"korin/pkg/klabels"
	"strings"
)

func Return(line string) klabels.Label {
	var values []string
	var prevChar rune

	inQuotations := false
	currentToken := ""

	for i, char := range line[7:] {
		if char == ',' && !inQuotations {
			values = append(values, strings.TrimSpace(currentToken))
			currentToken = ""
			prevChar = char
			continue
		}
		if prevChar != '\\' && char == '"' {
			inQuotations = !inQuotations
		}
		currentToken += string(char)

		if prevChar == '/' && char == '/' {
			values = append(values, strings.TrimSpace(currentToken))
			currentToken = ""
			prevChar = char
			break
		}

		if i == len(line)-8 {
			values = append(values, strings.TrimSpace(currentToken))
			currentToken = ""
			prevChar = char
			break
		}

		prevChar = char
	}
	return klabels.Label{Kind: klabels.ReturnKind, Data: klabels.ReturnStatement{Values: values}}
}
