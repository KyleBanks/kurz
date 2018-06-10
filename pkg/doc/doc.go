package doc

import (
	"io"
)

type Style int

const (
	Normal Style = 0

	Bold Style = 10 << iota
	Italic
	Underline

	BlockQuote Style = 100 << iota
	Code
	CodeBlock
	Image
	Link
	Unknown
)

type Styler interface {
	Style(string, Style) string
}

type Resolver interface {
	Resolve(string) (io.ReadCloser, error)
}

type Parser interface {
	Parse(io.Reader) (Document, error)
}

type Document struct {
	Headings []Heading
}

type Heading struct {
	Title   string
	Level   int
	Content []Section
}

type Section struct {
	Text string
}

func NewDocument(path string, r Resolver, p Parser) (Document, error) {
	content, err := r.Resolve(path)
	if err != nil {
		return Document{}, err
	}
	defer content.Close()

	return p.Parse(content)
}

// NopStyler implements a no-op Styler.
type NopStyler struct{}

// Style returns the provided string with no styling applied.
func (NopStyler) Style(s string, st Style) string {
	return s
}
