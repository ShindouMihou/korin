package korin

import "strings"

type Writer struct {
	content string
}

func (writer *Writer) NextLine() {
	writer.content += "\n"
}

func (writer *Writer) WriteLine(line string) {
	writer.NextLine()
	writer.Write(line)
}

func (writer *Writer) Write(line string) {
	writer.content += line
}

func (writer *Writer) Has(content string) bool {
	return strings.Contains(writer.content, content)
}

func (writer *Writer) Contents() string {
	return writer.content
}
