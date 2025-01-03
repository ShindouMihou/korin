package kplugins

import (
	"github.com/ShindouMihou/korin/pkg/klabels"
)

type Plugin interface {
	Group() string
	Version() string
	Name() string
	Process(line string, index int, headers *Headers, stack []klabels.Analysis) (string, error)
}
