package labelers

import (
	"korin/pkg/klabels"
	"strings"
)

func Package(line string) klabels.Label {
	name := strings.SplitN(line[8:], " ", 2)[0] // "package " is 8 chars
	return klabels.Label{Kind: klabels.PackageKind, Data: klabels.PackageDeclaration{Name: name}}
}
