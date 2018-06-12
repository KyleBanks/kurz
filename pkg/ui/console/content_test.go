package console

import "testing"

func TestContentState_get(t *testing.T) {
	c := newContentState()
	var state string

	// Un-set
	state = c.get(1, 1)
	if state != "" {
		t.Errorf("Unexpected result for un-set state, expected='', got=%v", state)
	}

	// Set
	c.store["1:1"] = "STATE"
	state = c.get(1, 1)
	if state != "STATE" {
		t.Errorf("Unexpected result for set state, expected=STATE, got=%v", state)
	}

	// Get other
	state = c.get(1, 2)
	if state != "" {
		t.Errorf("Unexpected result for other index, expected='', got=%v", state)
	}
}

func TestContentState_set(t *testing.T) {
	c := newContentState()
	var state string

	c.set(1, 1, "STATE")
	state = c.store["1:1"]
	if state != "STATE" {
		t.Errorf("Unexpected state, expected=STATE, got=%v", state)
	}

	// Set other
	c.set(1, 2, "OTHER")
	state = c.store["1:1"]
	if state != "STATE" {
		t.Errorf("Unexpected state after setting other, expected=STATE, got=%v", state)
	}
}

func TestContentState_clear(t *testing.T) {
	c := newContentState()
	var state string

	// Clear
	c.store["1:1"] = "STATE"
	c.clear(1, 1)
	state = c.store["1:1"]
	if state != "" {
		t.Errorf("Unexpected state after clear, expected='', got=%v", state)
	}

	// Clear other
	c.store["1:1"] = "STATE"
	c.clear(1, 2)
	state = c.store["1:1"]
	if state != "STATE" {
		t.Errorf("Unexpected state after clearing other, expected=STATE, got=%v", state)
	}
}

func TestContentKey(t *testing.T) {
	k := contentKey(1, 4)
	if k != "1:4" {
		t.Fatalf("Unexpected contentKey, expected=1:4, got=%v", k)
	}
}
