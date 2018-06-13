package console

import "testing"

func TestInput_String(t *testing.T) {
	i := input{
		symbol: "SYMBOL",
		label:  "LABEL",
	}

	expect := "[black:white:b]SYMBOL[-:-:-] LABEL"
	got := i.String()
	if expect != got {
		t.Errorf("Unexpected result, expected=%v, got=%v", expect, got)
	}
}

func TestNewInputHandler(t *testing.T) {
	w := NewWindow()
	i := newInputHandler(w)

	if i.w != w {
		t.Fatalf("Unexpected window, expected=%v, got=%v", w, i.w)
	}

	if w.GetInputCapture() == nil {
		t.Fatal("Unexpected nil input capture")
	}

	if i.tableOfContents == nil || i.content == nil {
		t.Fatal("Expected inputs to be set")
	}
}
