package debug

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var Enabled = false
var Out io.Writer = os.Stdout

// Log outputs a debug log message (with formatting) if
// the Enabled flag is set to true.
func Log(msg string, a ...interface{}) {
	if !Enabled {
		return
	}

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(Out, msg, a...)
}
