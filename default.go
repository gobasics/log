package log

import "os"

const (
	Skip       = 2
	TimeFormat = "Mon, 02 Jan 2006 15:04:05.999 UTC"
	Verbosity  = 255
)

var Writer = os.Stdout

var V = New()

var Fatal = V(0)

var Error = V(1)

var Warning = V(2)

var Info = V(3)
