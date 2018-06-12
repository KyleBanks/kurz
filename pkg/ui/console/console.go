package console

import (
	"bytes"
	"fmt"

	"github.com/KyleBanks/kurz/pkg/debug"
	"github.com/KyleBanks/kurz/pkg/doc"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type FocusMode int

const (
	FocusTableOfContents FocusMode = iota
	FocusContent
)

type Window struct {
	*tview.Application

	root *tview.Flex

	modal *tview.Modal

	tableOfContents *tview.List
	contentBody     *tview.TextView
	commandBar      *tview.TextView

	doc doc.Document

	focusMode       FocusMode
	selectedHeader  int
	selectedSection int

	contentState *contentState
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
		AddItem(w.CommandBar(), 1, 1, false)

	w.Application = tview.NewApplication().
		SetRoot(w.root, true)

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
	w.setFocusMode(FocusTableOfContents)
}

func (w *Window) setFocusMode(f FocusMode) {
	w.focusMode = f

	switch f {
	case FocusTableOfContents:
		w.SetInputCapture(w.tableOfContentsInputHandler)
		w.SetFocus(w.tableOfContents)
		w.ContentBody().Highlight()
	case FocusContent:
		w.SetInputCapture(w.contentInputHandler)
		w.SetFocus(w.ContentBody())
		w.setSelectedSection(0)
	}

	w.renderCommandBar()

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

func (w *Window) CommandBar() *tview.TextView {
	if w.commandBar == nil {
		w.commandBar = tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(false)
	}

	return w.commandBar
}

func (w *Window) renderCommandBar() {
	w.commandBar.Clear()

	var cmds []command
	switch w.focusMode {

	case FocusTableOfContents:
		cmds = commandsTableOfContents
	case FocusContent:
		cmds = commandsContent
	}

	var buf bytes.Buffer
	for _, c := range cmds {
		buf.WriteString(c.String())
		buf.WriteString("   ")
	}
	w.commandBar.SetText(buf.String())
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
	if idx < 0 || idx >= len(w.getSelectedHeader().Content) {
		return
	}

	if state := w.contentState.get(w.selectedHeader, idx); state == collapsedContent {
		w.contentState.clear(w.selectedHeader, idx)
	} else {
		w.contentState.set(w.selectedHeader, idx, collapsedContent)
	}

	w.renderContentBody()
}

func (w *Window) tableOfContentsInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {

	// Exit the application
	case tcell.KeyEscape:
		w.Stop()

	// Focus on the content
	case tcell.KeyRight:
		fallthrough
	case tcell.KeyEnter:
		w.setFocusMode(FocusContent)
		return nil

	// Ignore keys, don't bubble up the event.
	case tcell.KeyLeft:
		return nil
	}

	return event
}

func (w *Window) contentInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {

	// Go back to the Table of Contents
	case tcell.KeyLeft:
		fallthrough
	case tcell.KeyEscape:
		w.setFocusMode(FocusTableOfContents)
		return nil

	// Manage content selection
	case tcell.KeyUp:
		w.setSelectedSection(w.selectedSection - 1)
	case tcell.KeyDown:
		w.setSelectedSection(w.selectedSection + 1)

	case tcell.Key(256):
		w.collapseSection(w.selectedSection)

	// Ignore keys, don't bubble up the event.
	case tcell.KeyRight:
		return nil
	}

	return event
}

func (w *Window) tableOfContentsSelectionHandler(index int, mainText, secondaryText string, shortcut rune) {
	w.setSelectedHeader(index)
}
