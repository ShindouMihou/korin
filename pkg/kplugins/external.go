package kplugins

import (
	"github.com/ShindouMihou/korin/pkg/klabels"
)

type Plugin interface {
	Group() string
	Version() string
	Name() string
	Context(file string) *any
	FreeContext(file string)
	Process(line string, index int, headers *Headers, stack []klabels.Analysis, context *any) (string, error)
}

const SkipLine = "{$(SKIP_LINE)$+KORIN"
const NoChanges = ""
