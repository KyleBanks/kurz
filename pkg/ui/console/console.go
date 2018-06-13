package console

import (
	"bytes"
	"fmt"

	"github.com/KyleBanks/kurz/pkg/debug"
	"github.com/KyleBanks/kurz/pkg/doc"

	"github.com/atotto/clipboard"
	"github.com/rivo/tview"
)

type focusMode int

const (
	focusTableOfContents focusMode = iota
	focusContent
)

type Window struct {
	*tview.Application

	root *tview.Flex

	modal *tview.Modal

	tableOfContents *tview.List
	contentBody     *tview.TextView
	inputBar        *tview.TextView

	doc doc.Document

	focusMode       focusMode
	selectedHeader  int
	selectedSection int

	contentState *contentState
	inputHandler *inputHandler
}

func NewWindow() *Window {
	w := Window{
		modal:        tview.NewModal(),
		contentState: newContentState(),
	}

	w.root = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewGrid().
			SetBorders(true).
			AddItem(w.TableOfContents(), 0, 0, 1, 1, 0, 0, false).
			AddItem(w.ContentBody(), 0, 1, 1, 3, 0, 0, false),
			0, 1, false).
		AddItem(w.InputBar(), 1, 1, false)

	w.Application = tview.NewApplication().
		SetRoot(w.root, true)

	w.inputHandler = newInputHandler(&w)

	return &w
}

func (w *Window) ShowMessage(msg string) {
	w.modal.SetText(msg)
	w.SetRoot(w.modal, true)
}

func (w *Window) HideMessage() {
	w.SetRoot(w.root, true)
}

func (w *Window) RenderDocument(d doc.Document) {
	w.HideMessage()

	w.doc = d
	w.renderTableOfContents()
	w.setFocusMode(focusTableOfContents)
}

func (w *Window) setFocusMode(f focusMode) {
	w.focusMode = f

	switch f {
	case focusTableOfContents:
		w.SetFocus(w.tableOfContents)
		w.ContentBody().Highlight()
	case focusContent:
		w.SetFocus(w.ContentBody())
		w.setSelectedSection(0)
	}

	w.inputHandler.setFocusMode(f)
	w.renderInputBar()

	// Draw to ensure that the highlight is updated. This is important when
	// switching focus as the highlight functions don't always trigger a
	// new draw.
	w.Draw()
}

func (w *Window) TableOfContents() *tview.List {
	if w.tableOfContents == nil {
		w.tableOfContents = tview.NewList().
			SetChangedFunc(w.tableOfContentsSelectionHandler)
	}

	return w.tableOfContents
}

func (w *Window) renderTableOfContents() {
	w.tableOfContents.Clear()
	for _, h := range w.doc.Headers {
		text := h.Title
		if h.Level <= 2 {
			text = Styler{}.Style(text, doc.Bold)
		}

		for i := 0; i < h.Level-1; i++ {
			text = " " + text
		}

		w.tableOfContents.AddItem(text, "", 0, nil)
	}
}

func (w *Window) ContentBody() *tview.TextView {
	if w.contentBody == nil {
		w.contentBody = tview.NewTextView().
			SetDynamicColors(!debug.Enabled).
			SetRegions(true).
			SetScrollable(true).
			SetWrap(true).
			SetWordWrap(true)
	}

	return w.contentBody
}

func (w *Window) renderContentBody() {
	w.contentBody.Clear()
	if w.doc.Headers == nil {
		return
	}

	var buf bytes.Buffer
	for i, s := range w.getSelectedHeader().Content {
		text := s.Text
		if special := w.contentState.get(w.selectedHeader, i); len(special) > 0 {
			text = special
		}

		buf.WriteString(fmt.Sprintf(`["%d"]%v[""]`, i, text))
		buf.WriteString("\n")
	}
	w.contentBody.SetText(buf.String())
	w.contentBody.ScrollToBeginning()
}

func (w *Window) InputBar() *tview.TextView {
	if w.inputBar == nil {
		w.inputBar = tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(false)
	}

	return w.inputBar
}

func (w *Window) renderInputBar() {
	w.inputBar.Clear()
	w.inputBar.SetText(w.inputHandler.String())
}

func (w *Window) getSelectedHeader() doc.Header {
	return w.doc.Headers[w.selectedHeader]
}

func (w *Window) setSelectedHeader(selected int) {
	if selected < 0 || selected >= len(w.doc.Headers) {
		return
	}

	w.selectedHeader = selected
	w.renderContentBody()
}

func (w *Window) setSelectedSection(selected int) {
	numSections := len(w.getSelectedHeader().Content)
	if selected < 0 {
		selected = numSections - 1
	} else if selected >= numSections {
		selected = 0
	}

	w.selectedSection = selected
	w.contentBody.Highlight(fmt.Sprintf("%d", selected))
	w.contentBody.ScrollToHighlight()
}

func (w *Window) collapseSection(idx int) {
	if !w.isValidSectionIndex(idx) {
		return
	}

	if state := w.contentState.get(w.selectedHeader, idx); state == collapsedContent {
		w.contentState.clear(w.selectedHeader, idx)
	} else {
		w.contentState.set(w.selectedHeader, idx, collapsedContent)
	}

	w.renderContentBody()
}

func (w *Window) copySection(idx int) {
	if !w.isValidSectionIndex(idx) {
		return
	}

	// Use GetRegionText to have the formatting stripped from the
	// content, including colors/bolding/etc.
	text := w.contentBody.GetRegionText(fmt.Sprintf("%d", idx))
	clipboard.WriteAll(text)
}

func (w *Window) isValidSectionIndex(idx int) bool {
	return idx >= 0 && idx < len(w.getSelectedHeader().Content)
}

func (w *Window) tableOfContentsSelectionHandler(index int, mainText, secondaryText string, shortcut rune) {
	w.setSelectedHeader(index)
}
