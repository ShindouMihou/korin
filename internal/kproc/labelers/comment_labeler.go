package labelers

import (
	"korin/pkg/klabels"
	"strings"
)

func Comment(line string) *klabels.Label {
	var prevToken rune
	startingIndex := -1
	for index, char := range line {
		if char == '/' && prevToken == '/' {
			startingIndex = index + 1
			break
		}
		prevToken = char
	}

	if startingIndex != -1 {
		return &klabels.Label{Kind: klabels.CommentKind, Data: strings.TrimSpace(line[startingIndex:])}
	}
	return nil
}
