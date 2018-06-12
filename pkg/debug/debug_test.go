package debug

import (
	"bytes"
	"testing"
)

func TestLog(t *testing.T) {

	// Disabled
	output := bytes.Buffer{}
	Out = &output
	Enabled = false
	Log("Message")
	if len(output.String()) > 0 {
		t.Errorf("Unexpected logging when disabled, expected='', got=%v", output.String())
	}

	// Enabled
	output = bytes.Buffer{}
	Out = &output
	Enabled = true
	Log("Message")
	if output.String() != "Message\n" {
		t.Errorf("Unexpected log, expected=Message\n, got=%v", output.String())
	}
}
