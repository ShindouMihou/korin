package korin

import (
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/korin/pkg/klabels"
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

func (assistant Reader) Annotations(prefix string, annotations string) []string {
	var foundAnnotations []string
	for _, annotation := range strings.Split(annotations, " ") {
		if kstrings.HasPrefix(annotation, prefix+":\"") {
			currentAnnotations := strings.Split(annotation[len(prefix)+2:len(annotation)-1], ",")
			foundAnnotations = append(foundAnnotations, currentAnnotations...)
		}
	}
	return foundAnnotations
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
