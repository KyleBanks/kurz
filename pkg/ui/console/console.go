package console

import (
	"bytes"
	"fmt"

	"github.com/KyleBanks/kurz/pkg/doc"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type FocusMode int

const (
	FocusTableOfContents FocusMode = iota
	FocusContent         FocusMode = iota
)

type Window struct {
	*tview.Application

	root *tview.Grid

	modal *tview.Modal

	tableOfContents     tview.Primitive
	tableOfContentsList *tview.List
	contentBody         *tview.TextView

	doc doc.Document

	focusMode       FocusMode
	selectedHeading int
	selectedSection int
}

func NewWindow() *Window {
	var w Window

	w.modal = tview.NewModal()

	w.root = tview.NewGrid().
		SetBorders(true).
		AddItem(w.TableOfContents(), 0, 0, 1, 1, 0, 0, false).
		AddItem(w.ContentBody(), 0, 1, 1, 3, 0, 0, false)

	w.Application = tview.NewApplication().
		SetRoot(w.root, true)

	return &w
}

func (w *Window) ShowMessage(msg string) {
	w.modal.SetText(msg)
	w.SetRoot(w.modal, false)
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
		w.SetFocus(w.tableOfContentsList)
	case FocusContent:
		w.SetInputCapture(w.contentInputHandler)
		w.SetFocus(w.ContentBody())
	}
}

func (w *Window) TableOfContents() tview.Primitive {
	if w.tableOfContents == nil {
		header := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetDynamicColors(true).
			SetText("[purple::bu]Table of Contents")

		w.tableOfContentsList = tview.NewList().
			SetChangedFunc(w.tableOfContentsSelectionHandler)

		w.tableOfContents = tview.NewGrid().
			AddItem(header, 0, 0, 1, 1, 0, 0, false).
			AddItem(w.tableOfContentsList, 1, 0, 10, 1, 0, 0, false)
	}

	return w.tableOfContents
}

func (w *Window) renderTableOfContents() {
	w.tableOfContentsList.Clear()
	for _, h := range w.doc.Headings {
		// TODO: Use first char as the shortcut
		w.tableOfContentsList.AddItem(fmt.Sprintf("%-20s", h.Title), "", 0, nil)
	}
}

func (w *Window) ContentBody() *tview.TextView {
	if w.contentBody == nil {
		w.contentBody = tview.NewTextView().
			SetRegions(true).
			SetWrap(true).
			SetWordWrap(true).
			SetScrollable(true)
	}

	return w.contentBody
}

func (w *Window) renderContentBody() {
	w.contentBody.Clear()
	if w.doc.Headings == nil {
		return
	}

	var buf bytes.Buffer
	for i, s := range w.doc.Headings[w.selectedHeading].Content {
		buf.WriteString(fmt.Sprintf(`["%d"]%v[""]`, i, s.Text))
	}
	w.contentBody.SetText(buf.String())
	w.setSelectedSection(0)
}

func (w *Window) setSelectedSection(selected int) {
	numSections := len(w.doc.Headings[w.selectedHeading].Content)
	if selected < 0 {
		selected = 0
	} else if selected >= numSections {
		selected = numSections - 1
	}

	w.selectedSection = selected
	w.contentBody.Highlight(fmt.Sprintf("%d", selected))
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

	// Ignore keys, don't bubble up the event.
	case tcell.KeyRight:
		return nil
	}

	return event
}

func (w *Window) tableOfContentsSelectionHandler(index int, mainText, secondaryText string, shortcut rune) {
	w.selectedHeading = index
	w.renderContentBody()
}
