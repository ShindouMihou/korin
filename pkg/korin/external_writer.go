package korin

import "strings"

type Writer struct {
	content string
}

// NextLine adds a new line to the content.
func (writer *Writer) NextLine() {
	writer.content += "\n"
}

// WriteLine adds a new line and writes the line.
func (writer *Writer) WriteLine(line string) {
	writer.NextLine()
	writer.Write(line)
}

// Write writes a line to the content.
func (writer *Writer) Write(line string) {
	writer.content += line
}

// Has checks if the content contains the given string.
func (writer *Writer) Has(content string) bool {
	return strings.Contains(writer.content, content)
}

// Contents returns the content of the writer.
func (writer *Writer) Contents() string {
	return writer.content
}
