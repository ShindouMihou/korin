package kplugins

import "github.com/ShindouMihou/korin/pkg/klabels"

type AnalysisHelperType int

var AnalysisHelper AnalysisHelperType = 0

func (a AnalysisHelperType) HasOpenBracket(labels []klabels.Label) bool {
	return ReadHelper.Get(klabels.ScopeBeginKind, labels) != nil
}

func (a AnalysisHelperType) HasClosingBracket(labels []klabels.Label) bool {
	return ReadHelper.Get(klabels.ScopeEndKind, labels) != nil
}

// CheckMultiLineConstOrVar checks if the current const or var is multi-lined, this does it by checking if the
// lines before the current line has a ScopeBegin label, and if it has a ScopeEnd label before a ScopeBegin, it will
// consider it as a non-multi-lined const or var.
func (a AnalysisHelperType) CheckMultiLineConstOrVar(currentIndex int, stack []klabels.Analysis) bool {
	analysis := stack[currentIndex]
	hasScopeBegin, hasScopeEnd := false, false
	for index := analysis.Line; hasScopeBegin == false && index > 0; index-- {
		for _, label := range stack[index].Labels {
			label := label
			if label.Kind == klabels.ConstScopeBeginKind || label.Kind == klabels.VarScopeBeginKind {
				hasScopeBegin = true
			}
			if label.Kind == klabels.ConstScopeEndKind || label.Kind == klabels.VarScopeEndKind {
				hasScopeEnd = true
				if !hasScopeBegin {
					break
				}
			}
			if !hasScopeBegin && hasScopeEnd {
				break
			}
		}
		if !hasScopeBegin && hasScopeEnd {
			break
		}
	}
	return hasScopeBegin
}
