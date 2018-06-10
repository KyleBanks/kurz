package parser

import (
	"bytes"
	"testing"

	"github.com/KyleBanks/kurz/pkg/doc"
)

func TestMarkdown_Parse_Basic(t *testing.T) {
	m := NewMarkdown(doc.NopStyler{})

	d, err := m.Parse(bytes.NewBufferString(`
# Heading 1

Text goes here.

Another paragraph.

## Heading 2

And here's some ` + "`" + `code` + "`" + `.
	`))
	if err != nil {
		t.Fatal(err)
	}

	assertDocsEqual(t, d, doc.Document{
		Headings: []doc.Heading{
			{Title: "Heading 1", Level: 1, Content: []doc.Section{
				{"Text goes here.\n"},
				{"Another paragraph.\n"},
			}},
			{Title: "Heading 2", Level: 2, Content: []doc.Section{
				{"And here's some code.\n"},
			}},
		},
	})
}

func assertDocsEqual(t *testing.T, got, exp doc.Document) {
	if len(got.Headings) != len(exp.Headings) {
		t.Fatalf("Heading count mismatch, expected=%v, got=%v", len(exp.Headings), len(got.Headings))
	}

	for h := 0; h < len(got.Headings); h++ {
		h1 := got.Headings[h]
		h2 := exp.Headings[h]
		if h1.Title != h2.Title {
			t.Errorf("[h=%v] Unexpected title, expected=%v, got=%v", h, h2.Title, h1.Title)
		}
		if h1.Level != h2.Level {
			t.Errorf("[h=%v] Unexpected level, expected=%v, got=%v", h, h2.Level, h1.Level)
		}

		if len(h1.Content) != len(h2.Content) {
			t.Fatalf("[h=%v] Content length mismatch, expected=%v, got=%v", h, len(h2.Content), len(h1.Content))
		}

		for c := 0; c < len(h1.Content); c++ {
			c1 := h1.Content[c]
			c2 := h2.Content[c]

			if c1.Text != c2.Text {
				t.Errorf("[h=%v, c=%v] Unexpected text, expected'=%v', got'=%v'", h, c, c2.Text, c1.Text)
			}
		}
	}
}
