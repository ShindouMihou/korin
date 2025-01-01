package labelers

import "github.com/ShindouMihou/korin/pkg/klabels"

func TypeDeclaration(line string) klabels.Label {
	label := klabels.Label{Kind: klabels.TypeDeclarationKind}

	name := ""
	kind := ""

	currentToken := ""
	currentLocation := "name"

	// type Name struct

	for _, char := range line[5:] {
		if char == ' ' {
			if currentLocation == "name" && currentToken != "" {
				name = currentToken

				currentToken = ""
				currentLocation = "kind"
				continue
			}
			if currentLocation == "kind" && currentToken != "" {
				kind = currentToken
				break
			}
			continue
		}
		currentToken += string(char)
	}

	label.Data = klabels.TypeDeclaration{Name: name, Kind: kind}
	return label
}
