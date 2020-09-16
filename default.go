package log

import "os"

const (
	TimeFormat = "Mon, 02 Jan 2006 15:04:05.999 UTC"
)

var V = New(
	WithSkip(3),
	WithTimeFormat(TimeFormat),
	WithVerbosity(255),
	WithWriter(os.Stdout),
)

var Fatal = V(0)

var Error = V(1)

var Warning = V(2)

var Info = V(3)
