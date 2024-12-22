package korin

import "slices"

type Headers struct {
	imports []string
	pkg     string
}

func (headers *Headers) Import(pkg string) {
	if slices.Contains(headers.imports, pkg) {
		return
	}
	headers.imports = append(headers.imports, pkg)
}

func (headers *Headers) Package(name string) {
	headers.pkg = name
}

func (headers *Headers) Format() string {
	text := "package " + headers.pkg
	text += "\n"
	text += "\n"
	if len(headers.imports) > 0 {
		text += WriteAssistant.Import(headers.imports)
		text += "\n"
	}
	return text
}
