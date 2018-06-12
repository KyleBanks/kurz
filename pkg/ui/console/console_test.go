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
		t.Error("Unexpected nil root")
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
	if w.contentState == nil {
		t.Errorf("Unexpected nil contentState")
	}
}

func TestWindow_RenderDocument(t *testing.T) {
	w := NewWindow()
	d := doc.Document{
		Headers: []doc.Header{
			{Title: "Header 1"},
			{Title: "Header 2"},
		},
	}

	w.RenderDocument(d)
	if !reflect.DeepEqual(d, w.doc) {
		t.Errorf("Unexpected w.doc, expected=%v, got=%v", d, w.doc)
	}

	if w.focusMode != FocusTableOfContents {
		t.Errorf("Unexpected w.focusMode, expected=%v, got=%v", w.focusMode, FocusTableOfContents)
	}

	if w.tableOfContents.GetItemCount() != len(d.Headers) {
		t.Errorf("Unexpected tableOfContents length, expected=%v, got=%v", len(d.Headers), w.tableOfContents.GetItemCount())
	}

	if w.selectedHeader != 0 {
		t.Errorf("Unexpected selectedHeader, expected=0, got=%v", w.selectedHeader)
	}
}

func TestWindow_getSelectedHeader(t *testing.T) {
	var w Window
	w.doc = doc.Document{
		Headers: []doc.Header{
			{Title: "Header 1"},
			{Title: "Header 2"},
		},
	}

	w.selectedHeader = 0
	got := w.getSelectedHeader().Title
	if w.getSelectedHeader().Title != "Header 1" {
		t.Errorf("Unexpected Title, expected=Header 1, got=%v", got)
	}

	w.selectedHeader = 1
	got = w.getSelectedHeader().Title
	if w.getSelectedHeader().Title != "Header 2" {
		t.Errorf("Unexpected Title, expected=Header 2, got=%v", got)
	}
}

func TestWindow_setSelectedHeader(t *testing.T) {
	w := NewWindow()
	w.doc = doc.Document{
		Headers: []doc.Header{
			{Title: "Header 1", Content: []doc.Section{
				{Text: "H1T1"},
				{Text: "H1T2"},
			}},
			{Title: "Header 2", Content: []doc.Section{
				{Text: "H2T1"},
				{Text: "H2T2"},
			}},
		},
	}

	tests := []struct {
		header int
		region string
		expect string
	}{
		{header: 0, region: "0", expect: "H1T1"},
		{header: 0, region: "1", expect: "H1T2"},
		{header: 1, region: "0", expect: "H2T1"},
		{header: 1, region: "1", expect: "H2T2"},
	}

	for idx, tt := range tests {
		w.setSelectedHeader(tt.header)
		if w.selectedHeader != tt.header {
			t.Errorf("[%d] Unexpected selectedHeader, expected=%d, got=%v", idx, tt.header, w.selectedHeader)
		}

		r1 := w.contentBody.GetRegionText(tt.region)
		if r1 != tt.expect {
			t.Errorf("[%d] Unexpected text at region %v, expected=%v, got=%v", idx, tt.region, tt.expect, r1)
		}
	}

	// Ignore invalid input
	w.setSelectedHeader(1)

	w.setSelectedHeader(10)
	if w.selectedHeader != 1 {
		t.Errorf("Unexpected selectedHeader for high input, expected=1, got=%v", w.selectedHeader)
	}

	w.setSelectedHeader(-1)
	if w.selectedHeader != 1 {
		t.Errorf("Unexpected selectedHeader for low input, expected=1, got=%v", w.selectedHeader)
	}
}
