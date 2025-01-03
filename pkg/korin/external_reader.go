package korin

import (
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/korin/pkg/klabels"
	"strings"
)

type ReadHelperType int

var ReadHelper ReadHelperType = 0

// Parameters is a helper function that gets the parameters based on the prefix and returns them to the caller.
func (assistant ReadHelperType) Parameters(prefix string, labels []klabels.Label) (bool, []string) {
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

// Annotations is a helper function that gets the annotations based on the prefix and returns them to the caller.
func (assistant ReadHelperType) Annotations(prefix string, annotations string) []string {
	var foundAnnotations []string
	for _, annotation := range strings.Split(annotations, " ") {
		if kstrings.HasPrefix(annotation, prefix+":\"") {
			currentAnnotations := strings.Split(annotation[len(prefix)+2:len(annotation)-1], ",")
			foundAnnotations = append(foundAnnotations, currentAnnotations...)
		}
	}
	return foundAnnotations
}

// Get is a helper function that gets the label based on the kind and returns it to the caller.
func (assistant ReadHelperType) Get(kind klabels.LabelKind, labels []klabels.Label) *klabels.Label {
	for _, label := range labels {
		label := label
		if label.Kind == kind {
			return &label
		}
	}

	return nil
}

// Filter is a helper function that filters the labels based on the kinds and returns the result to the caller.
func (assistant ReadHelperType) Filter(labels []klabels.Label, on func(label klabels.Label), kinds ...klabels.LabelKind) {
	for _, kind := range kinds {
		for _, label := range labels {
			label := label
			if label.Kind == kind {
				on(label)
			}
		}
	}
}

// PluginResultFunc is a function that returns a string and an error.
type PluginResultFunc[T any] func(res T) (string, error)

// Require is a helper function that requires a specific label kind and returns the result to the caller.
func (assistant ReadHelperType) Require(kind klabels.LabelKind, labels []klabels.Label, on PluginResultFunc[klabels.Label]) (string, error) {
	label := assistant.Get(kind, labels)
	if label != nil {
		return on(*label)
	}
	return "", nil
}

// GetParameters is a helper function that gets the parameters from the labels and returns them to the caller.
func (assistant ReadHelperType) GetParameters(prefix string, labels []klabels.Label, on PluginResultFunc[[]string]) (string, error) {
	shouldAnnotate, parameters := assistant.Parameters(prefix, labels)
	if shouldAnnotate {
		return on(parameters)
	}
	return "", nil
}
