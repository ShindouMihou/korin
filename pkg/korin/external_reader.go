package korin

import (
	"korin/internal/kstrings"
	"korin/pkg/klabels"
	"strings"
)

type Reader int

var ReadAssistant Reader = 0

func (assistant Reader) Parameters(prefix string, labels []klabels.Label) (bool, []string) {
	var parameters []string
	label := assistant.Get(klabels.CommentKind, labels)
	if label == nil {
		return false, parameters
	}

	value := label.Data.(string)
	if kstrings.HasPrefix(value, prefix) {
		if len(value) > len(prefix)+1 {
			for _, param := range strings.Split(value[len(prefix)+1:len(value)-1], ",") {
				if param == "" {
					continue
				}
				parameters = append(parameters, strings.TrimSpace(param))
			}
		}
		return true, parameters
	}
	return false, parameters
}

func (assistant Reader) Get(kind klabels.LabelKind, labels []klabels.Label) *klabels.Label {
	for _, label := range labels {
		label := label
		if label.Kind == kind {
			return &label
		}
	}

	return nil
}
