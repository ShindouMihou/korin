package labelers

import (
	"github.com/ShindouMihou/korin/pkg/klabels"
	"strings"
)

func FieldDeclaration(line string) klabels.Label {
	label := klabels.Label{Kind: klabels.FieldDeclarationKind}

	name := ""
	kind := ""
	annotation := ""

	currentToken := ""
	currentLocation := "name"

	// Name string `json:"name"`
	for _, char := range strings.TrimSpace(line) {
		if char == ' ' {
			if currentLocation == "name" && currentToken != "" {
				name = currentToken

				currentToken = ""
				currentLocation = "kind"
				continue
			}
			if currentLocation == "kind" && currentToken != "" {
				kind = currentToken

				currentToken = ""
				currentLocation = "empty"
				continue
			}
			continue
		}
		if char == '`' {
			if currentLocation == "empty" {
				currentLocation = "annotation"
				continue
			}
			if currentLocation == "annotation" {
				annotation = currentToken
				break
			}
			continue
		}
		currentToken += string(char)
	}

	label.Data = klabels.FieldDeclaration{Name: name, Type: kind, Annotations: annotation}
	return label
}
