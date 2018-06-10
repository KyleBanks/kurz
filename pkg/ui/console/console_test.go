package console

import (
	"reflect"
	"testing"

	"github.com/KyleBanks/kurz/pkg/doc"
)

func TestNewWindow(t *testing.T) {
	w := NewWindow()
	if w == nil {
		t.Fatal("Unexpected nil Window")
	}

	if w.root == nil {
		t.Error("Unexpected nil roo")
	}
	if w.modal == nil {
		t.Error("Unexpected nil modal")
	}
	if w.tableOfContents == nil {
		t.Error("Unexpected nil tableOfContents")
	}
	if w.contentBody == nil {
		t.Error("Unexpected nil contentBody")
	}
	if w.commandBar == nil {
		t.Error("Unexpected nil commandBar")
	}
}

func TestWindow_RenderDocument(t *testing.T) {
	w := NewWindow()
	d := doc.Document{
		Headings: []doc.Heading{
			{Title: "Heading 1"},
			{Title: "Heading 2"},
		},
	}

	w.RenderDocument(d)
	if !reflect.DeepEqual(d, w.doc) {
		t.Errorf("Unexpected w.doc, expected=%v, got=%v", d, w.doc)
	}

	if w.focusMode != FocusTableOfContents {
		t.Errorf("Unexpected w.focusMode, expected=%v, got=%v", w.focusMode, FocusTableOfContents)
	}

	if w.tableOfContents.GetItemCount() != len(d.Headings) {
		t.Errorf("Unexpected tableOfContents length, expected=%v, got=%v", len(d.Headings), w.tableOfContents.GetItemCount())
	}
}
