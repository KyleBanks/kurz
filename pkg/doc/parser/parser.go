package parser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/KyleBanks/kurz/pkg/doc"

	"gopkg.in/russross/blackfriday.v2"
)

type Markdown struct{}

func (m Markdown) Parse(content io.ReadCloser) (doc.Document, error) {
	defer content.Close()

	b, err := ioutil.ReadAll(content)
	if err != nil {
		return doc.Document{}, err
	}

	md := blackfriday.New()
	root := md.Parse(b)

	var d doc.Document
	root.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.GoToNext
		}

		if node.Type != blackfriday.Heading || node.HeadingData.Level > 2 {
			return blackfriday.GoToNext
		}

		d.Headings = append(d.Headings, doc.Heading{
			Title:   m.nodeContents(node),
			Content: m.sectionContents(node),
		})

		return blackfriday.SkipChildren
	})

	return d, nil
}

func (m Markdown) sectionContents(heading *blackfriday.Node) []doc.Section {
	var sections []doc.Section

	n := heading.Next
	for n != nil && (n.Type != blackfriday.Heading || n.HeadingData.Level > 2) {
		sections = append(sections, m.newSection(n))
		n = n.Next
	}

	return sections
}

func (m Markdown) newSection(container *blackfriday.Node) doc.Section {
	var buf bytes.Buffer

	n := container
	for n != nil {
		str := strings.TrimSpace(m.nodeContents(n))
		if len(str) > 0 {
			buf.WriteString(str)
			if m.appendNewline(n) {
				buf.WriteString("\n")
			}
		}

		if n.FirstChild != nil && !m.skipChild(n) {
			n = n.FirstChild
		} else if n.Next != nil {
			n = n.Next
		} else if n.Parent != container {
			n = n.Parent.Next
		} else {
			n = nil
		}
	}

	text := buf.String()
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	return doc.Section{
		Text: text,
	}
}

func (Markdown) skipChild(n *blackfriday.Node) bool {
	switch n.Type {

	case blackfriday.Heading:
		fallthrough
	case blackfriday.Link:
		fallthrough
	case blackfriday.List:
		fallthrough
	case blackfriday.Image:
		return true

	default:
		return false
	}
}

func (Markdown) appendNewline(n *blackfriday.Node) bool {
	switch n.Type {

	case blackfriday.Code:
		return strings.Contains(string(n.Literal), "\n")

	case blackfriday.Paragraph:
		fallthrough
	case blackfriday.Heading:
		fallthrough
	case blackfriday.Image:
		return true

	default:
		return false
	}
}

func (m Markdown) nodeContents(n *blackfriday.Node) string {
	switch n.Type {

	case blackfriday.Heading:
		// The following is required when the root element of a header is a non-text type,
		// for instance a code block:
		// # `code`
		if n.FirstChild.Next != nil {
			return m.nodeContents(n.FirstChild.Next)
		}
		return m.nodeContents(n.FirstChild)

	// Text nodes
	case blackfriday.Code:
		return fmt.Sprintf("`%s`", n.Literal)
	case blackfriday.Item:
		fallthrough
	case blackfriday.Paragraph:
		fallthrough
	case blackfriday.Text:
		return string(n.Literal)

	// Skip nodes
	case blackfriday.Document:
		return ""

	// Special nodes
	case blackfriday.Link:
		text := m.nodeContents(n.FirstChild)
		if len(text) > 0 {
			return fmt.Sprintf("%v <%s>", text, n.LinkData.Destination)
		}
		return fmt.Sprintf("<%s>", n.LinkData.Destination)
	case blackfriday.List:
		item := n.FirstChild
		var buf bytes.Buffer
		for item != nil {
			buf.WriteString(fmt.Sprintf("%s %v\n", []byte{item.BulletChar}, strings.TrimSpace(m.nodeContents(item))))
			item = item.Next
		}
		return buf.String()

		// TODO: image

	default:
		return fmt.Sprintf("Type=%v\n", n.Type)
	}
}
