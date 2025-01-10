package kmod

import (
	"errors"
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/siopao/siopao"
	"strings"
)

type Module struct {
	Name string
}

func ReadModule(path string) (*Module, error) {
	if !kstrings.HasSuffix(path, "go.mod") {
		return nil, errors.New("path is not a go.mod file")
	}
	file := siopao.Open(path)

	reader, err := file.TextReader()
	if err != nil {
		return nil, err
	}

	module := Module{}
	if err := reader.EachLine(func(line string) {
		if kstrings.HasPrefix(line, "module ") {
			module.Name = strings.TrimPrefix(line, "module ")
		}
	}); err != nil {
		return nil, err
	}
	return &module, nil
}
