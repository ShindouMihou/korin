package kplugins

import (
	"github.com/ShindouMihou/korin/internal/kslices"
	"slices"
	"strings"
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
func (headers *Headers) Format(moduleName string, buildDir string) string {
	text := "package " + headers.pkg
	text += "\n"
	if len(headers.imports) > 0 {
		text += "\n"

		// Important part of headers modification since if we don't do this, then
		// all the source files will be imported from the original module path, which makes this preprocessing
		// entirely useless.
		moduleName = strings.TrimSuffix(moduleName, "/")
		buildDir = strings.TrimSuffix(buildDir, "/")
		for ind, imp := range headers.imports {
			if strings.HasPrefix(imp, moduleName+"/") {
				headers.imports[ind] = strings.ReplaceAll(imp, moduleName+"/", moduleName+"/"+buildDir+"/")
			}
		}

		text += SyntaxHelper.Import(headers.imports)
	}
	return text
}
