package kplugins

import (
	"github.com/ShindouMihou/korin/internal/kslices"
	"slices"
)

type Headers struct {
	imports []string
	pkg     string
}

// Import imports a package into the headers.
func (headers *Headers) Import(pkg string) {
	if slices.Contains(headers.imports, pkg) {
		return
	}
	headers.imports = append(headers.imports, pkg)
}

// RemoveImport removes a package from the headers.
func (headers *Headers) RemoveImport(pkg string) {
	headers.imports = kslices.RemoveString(headers.imports, pkg)
}

// Package sets the package name for the headers.
func (headers *Headers) Package(name string) {
	headers.pkg = name
}

// Format returns the formatted headers, this properly formats them, in which case, multiple imports will be
// multi-lined and the package will be at the top.
func (headers *Headers) Format() string {
	text := "package " + headers.pkg
	text += "\n"
	if len(headers.imports) > 0 {
		text += "\n"
		text += SyntaxHelper.Import(headers.imports)
	}
	return text
}
