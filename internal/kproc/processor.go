package kproc

import (
	"github.com/ShindouMihou/korin/internal/kproc/labelers"
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/korin/pkg/klabels"
	"strings"
)

func LabelLine(index int, line string) klabels.Analysis {
	var labels []klabels.Label

	line = strings.TrimSpace(line)
	if line == "(" || line == ")" || len(line) < 3 {
		return klabels.Analysis{Line: index}
	}

	if kstrings.HasSuffix(line, "{") {
		labels = append(labels, klabels.Label{Kind: klabels.ScopeBeginKind})
	}

	if kstrings.HasPrefix(line, "func") {
		labels = append(labels, labelers.FuncLine(line))
	} else if kstrings.HasPrefix(line, "const") {
		if !kstrings.HasSuffix(line, "(") {
			labels = append(labels, klabels.Label{Kind: klabels.ConstDeclarationKind})
			labels = append(labels, labelers.VariableAssignment(false, line[6:]))
		} else {
			labels = append(labels, klabels.Label{Kind: klabels.ConstScopeBeginKind})
		}
	} else if kstrings.HasPrefix(line, "type") {
		labels = append(labels, labelers.TypeDeclaration(line))
		for _, char := range line {
			if char == '{' {
				labels = append(labels, klabels.Label{Kind: klabels.ScopeBeginKind})
			}
			if char == '}' {
				labels = append(labels, klabels.Label{Kind: klabels.ScopeEndKind})
			}
		}
	} else if kstrings.HasPrefix(line, "var") {
		if !kstrings.HasSuffix(line, "(") {
			labels = append(labels, klabels.Label{Kind: klabels.VarDeclarationKind})
			labels = append(labels, labelers.VariableAssignment(false, line[4:]))
		} else {
			labels = append(labels, klabels.Label{Kind: klabels.VarScopeBeginKind})
		}
	} else if kstrings.HasPrefix(line, "package ") {
		labels = append(labels, labelers.Package(line))
	} else if kstrings.HasPrefix(line, "return ") {
		labels = append(labels, labelers.Return(line))
	} else {
		var previousToken rune
		isInsideExclusionZone := false

		for _, char := range line {
			if char == '=' && previousToken == ':' && !isInsideExclusionZone {
				// variable declaration
				labels = append(labels, labelers.VariableAssignment(false, line))
			} else if char == '=' && previousToken == ' ' && !isInsideExclusionZone {
				// variable reassignment
				labels = append(labels, labelers.VariableAssignment(true, line))
			} else if char == '"' && !isInsideExclusionZone {
				isInsideExclusionZone = true
			} else if char == '/' && previousToken == '/' && !isInsideExclusionZone {
				isInsideExclusionZone = true
			} else if char == '"' && isInsideExclusionZone {
				isInsideExclusionZone = false
			}
			previousToken = char
		}
	}

	if strings.Contains(line, "func(") {
		labels = append(labels, labelers.AnonymousFunction(line))
	}

	comment := labelers.Comment(line)
	if comment != nil {
		labels = append(labels, *comment)
	}

	if line == "}" {
		labels = append(labels, klabels.Label{Kind: klabels.ScopeEndKind})
	}

	return klabels.Analysis{Line: index, Labels: labels}
}
