package doc

import (
	"io"
)

type Resolver interface {
	Resolve(string) (io.ReadCloser, error)
}

type Parser interface {
	Parse(io.ReadCloser) (Document, error)
}

type Document struct {
	Headings []Heading
}

type Heading struct {
	Title   string
	Content []Section
}

type Section struct {
	Text string
	// TODO: type
}

func NewDocument(path string, r Resolver, p Parser) (Document, error) {
	content, err := r.Resolve(path)
	if err != nil {
		return Document{}, err
	}

	return p.Parse(content)
}
