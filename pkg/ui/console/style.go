package console

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/KyleBanks/kurz/pkg/doc"
)

var StyleMap = map[doc.Style]Style{
	doc.Bold:       Style{"", "", "b", ""},
	doc.Italic:     Style{"", "", "u", ""}, // There are no italics in a TUI.
	doc.Underline:  Style{"", "", "u", ""},
	doc.BlockQuote: Style{"", "", "b", "   "},
	doc.Code:       Style{"purple", "", "b", ""},
	doc.CodeBlock:  Style{"purple", "", "b", "   "},
	doc.Image:      Style{"#9331ee", "", "bu", ""},
	doc.Link:       Style{"green", "", "bu", ""},
	doc.Unknown:    Style{"red", "", "", ""},
}

var DefaultStyle Style

type Style struct {
	FgColor   string
	BgColor   string
	TextStyle string
	Indent    string
}

type Styler struct{}

func (Styler) Style(str string, ds doc.Style) string {
	s, ok := StyleMap[ds]
	if !ok {
		s = DefaultStyle
	}

	lines := strings.Split(str, "\n")
	var buf bytes.Buffer
	for i, l := range lines {
		buf.WriteString(s.Indent)
		buf.WriteString(fmt.Sprintf("[%v:%v:%v]", s.FgColor, s.BgColor, s.TextStyle))
		buf.WriteString(l)
		buf.WriteString("[-:-:-]")

		if i < len(lines)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
