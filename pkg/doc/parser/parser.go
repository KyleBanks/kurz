package parser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/KyleBanks/kurz/pkg/debug"
	"github.com/KyleBanks/kurz/pkg/doc"

	"gopkg.in/russross/blackfriday.v2"
)

type Markdown struct {
	Styler doc.Styler
}

func NewMarkdown(s doc.Styler) Markdown {
	return Markdown{
		Styler: s,
	}
}

func (m Markdown) Parse(r io.Reader) (doc.Document, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return doc.Document{}, err
	}

	md := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	root := md.Parse(b)

	var d doc.Document
	root.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.GoToNext
		} else if node.Type != blackfriday.Heading {
			return blackfriday.GoToNext
		}

		d.Headings = append(d.Headings, doc.Heading{
			Title:   m.nodeContents(node),
			Level:   node.HeadingData.Level,
			Content: m.sectionContents(node),
		})

		return blackfriday.SkipChildren
	})

	return d, nil
}

func (m Markdown) sectionContents(heading *blackfriday.Node) []doc.Section {
	var sections []doc.Section
	for n := heading.Next; n != nil && n.Type != blackfriday.Heading; n = n.Next {
		sections = append(sections, m.newSection(n))
	}
	return sections
}

func (m Markdown) newSection(container *blackfriday.Node) doc.Section {
	var buf bytes.Buffer
	var skipNext bool
	container.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			if !skipNext {
				str := m.nodeContents(n)
				buf.WriteString(str)
			}
		} else {
			if m.appendNewline(n) {
				buf.WriteString("\n")
			}
		}

		if m.skipChildren(n) {
			return blackfriday.SkipChildren
		}

		skipNext = m.skipNext(n)
		return blackfriday.GoToNext
	})

	text := buf.String()
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	return doc.Section{
		Text: text,
	}
}

func (Markdown) appendNewline(n *blackfriday.Node) bool {
	switch n.Type {

	case blackfriday.Code:
		fallthrough
	case blackfriday.CodeBlock:
		return !strings.HasSuffix(string(n.Literal), "\n")

	case blackfriday.Paragraph:
		return true

	default:
		return false
	}
}

// skipChildren returns true when the children of the provided Node
// should be skipped.
func (Markdown) skipChildren(n *blackfriday.Node) bool {
	switch n.Type {
	case blackfriday.Emph:
		fallthrough
	case blackfriday.Strong:
		return true
	}

	return false
}

// skipNext returns true when the next Node in the tree should be skipped.
func (Markdown) skipNext(n *blackfriday.Node) bool {
	switch n.Type {
	// Skip the next if the direct child is a text type, as it's
	// already been rendered and stylized. However, sometimes the child is
	// another type, such as an image in a link, in which case we want to render it.
	case blackfriday.Link:
		fallthrough
	case blackfriday.Image:
		if n.FirstChild != nil && n.FirstChild.Type == blackfriday.Text {
			return true
		}
	}

	return false
}

func (m Markdown) nodeContents(n *blackfriday.Node) string {
	if debug.Enabled {
		fmt.Printf("Type=%v, Literal=%s\n", n.Type, n.Literal)
	}

	switch n.Type {

	case blackfriday.BlockQuote:
		return m.Styler.Style(string(n.FirstChild.Literal), doc.BlockQuote)

	case blackfriday.Code:
		return m.Styler.Style(string(n.Literal), doc.Code)

	case blackfriday.CodeBlock:
		str := strings.TrimSpace(string(n.Literal))
		return m.Styler.Style(str, doc.CodeBlock)

	case blackfriday.Document:
		return ""

	case blackfriday.Emph:
		return m.Styler.Style(string(n.FirstChild.Literal), doc.Italic)

	case blackfriday.Item:
		return fmt.Sprintf("%s ", []byte{n.ListData.BulletChar})

	case blackfriday.List:
		return ""

	case blackfriday.Paragraph:
		return string(n.Literal)

	case blackfriday.Strong:
		return m.Styler.Style(string(n.FirstChild.Literal), doc.Bold)

	case blackfriday.Text:
		return string(n.Literal)

		// Special nodes

	case blackfriday.Heading:
		// The following is required when the root element of a header is a non-text type,
		// for instance a code block:
		// # `code`
		if n.FirstChild.Next != nil {
			return m.nodeContents(n.FirstChild.Next)
		}
		return m.nodeContents(n.FirstChild)

	case blackfriday.Link:
		text := m.nodeContents(n.FirstChild)
		if len(text) > 0 {
			text += " "
		}
		link := fmt.Sprintf("%v<%s> ", text, n.LinkData.Destination)
		return m.Styler.Style(link, doc.Link)

	case blackfriday.Image:
		text := m.nodeContents(n.FirstChild)
		if len(text) > 0 {
			text += " "
		}
		img := fmt.Sprintf("Image: %v<%s> ", text, n.LinkData.Destination)
		return m.Styler.Style(img, doc.Image)

	default:
		str := fmt.Sprintf("Unknown Node: {%v}", n)
		return m.Styler.Style(str, doc.Unknown)
	}
}
