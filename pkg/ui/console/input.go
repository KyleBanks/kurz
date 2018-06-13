package console

import (
	"bytes"
	"fmt"

	"github.com/gdamore/tcell"
)

type input struct {
	symbol  string
	label   string
	fn      func()
	swallow bool

	keys  []tcell.Key
	runes []rune
}

func (i input) String() string {
	if i.symbol == "" || i.label == "" {
		return ""
	}

	return fmt.Sprintf("[black:white:b]%v[-:-:-] %v", i.symbol, i.label)
}

type inputHandler struct {
	w     *Window
	focus focusMode

	tableOfContents []input
	content         []input
}

func newInputHandler(w *Window) *inputHandler {
	i := inputHandler{
		w: w,
	}

	w.SetInputCapture(i.handle)
	i.setInputs()

	return &i
}

func (i *inputHandler) handle(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyRune {
		return i.handleRune(e)
	} else {
		return i.handleKey(e)
	}
}

func (i *inputHandler) handleRune(e *tcell.EventKey) *tcell.EventKey {
	for _, input := range i.inputs() {
		for _, r := range input.runes {
			if r != e.Rune() {
				continue
			}

			if input.fn != nil {
				input.fn()
			}

			if input.swallow {
				return nil
			}
			return e
		}
	}
	return e
}

func (i *inputHandler) handleKey(e *tcell.EventKey) *tcell.EventKey {
	for _, input := range i.inputs() {
		for _, k := range input.keys {
			if k != e.Key() {
				continue
			}

			if input.fn != nil {
				input.fn()
			}

			if input.swallow {
				return nil
			}
			return e
		}
	}
	return e
}

func (i *inputHandler) setFocusMode(f focusMode) {
	i.focus = f
}

func (i *inputHandler) inputs() []input {
	var inputs []input
	switch i.focus {
	case focusTableOfContents:
		inputs = i.tableOfContents
	case focusContent:
		inputs = i.content
	}
	return inputs
}

func (i *inputHandler) String() string {
	var buf bytes.Buffer
	for _, c := range i.inputs() {
		str := c.String()
		if str == "" {
			continue
		}

		buf.WriteString(c.String())
		buf.WriteString("   ")
	}
	return buf.String()
}

func (i *inputHandler) setInputs() {
	i.tableOfContents = []input{
		{
			symbol: " ESC ",
			label:  "Exit",
			keys:   []tcell.Key{tcell.KeyEscape},
			fn:     i.w.Stop,
		},
		{
			symbol: "⬆ ",
			label:  "Up",
			keys:   []tcell.Key{tcell.KeyUp},
		},
		{
			symbol: "⬇ ",
			label:  "Down",
			keys:   []tcell.Key{tcell.KeyDown},
		},
		{
			symbol:  " ➡ / ENTER ",
			label:   "Select",
			keys:    []tcell.Key{tcell.KeyRight, tcell.KeyEnter},
			fn:      func() { i.w.setFocusMode(focusContent) },
			swallow: true,
		},
		{
			keys:    []tcell.Key{tcell.KeyLeft},
			swallow: true,
		},
	}

	i.content = []input{
		{
			symbol:  " ⬅ / ESC ",
			label:   "Go Back",
			keys:    []tcell.Key{tcell.KeyLeft, tcell.KeyEscape},
			fn:      func() { i.w.setFocusMode(focusTableOfContents) },
			swallow: true,
		},
		{
			symbol: "⬆ ",
			label:  "Up",
			keys:   []tcell.Key{tcell.KeyUp},
			fn:     func() { i.w.setSelectedSection(i.w.selectedSection - 1) },
		},
		{
			symbol: "⬇ ",
			label:  "Down",
			keys:   []tcell.Key{tcell.KeyDown},
			fn:     func() { i.w.setSelectedSection(i.w.selectedSection + 1) },
		},
		{
			symbol: " SPACE ",
			label:  "Collapse",
			runes:  []rune{32}, // space
			fn:     func() { i.w.collapseSection(i.w.selectedSection) },
		},
		{
			symbol:  " C ",
			label:   "Copy",
			runes:   []rune{99}, // c
			fn:      func() { i.w.copySection(i.w.selectedSection) },
			swallow: true,
		},
		{
			keys:    []tcell.Key{tcell.KeyRight},
			swallow: true,
		},
	}
}
