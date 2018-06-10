package debug

import "os"

var Enabled = os.Getenv("KURZ_DEBUG") == "true"
