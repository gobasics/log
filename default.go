package log

import (
	"os"
)

const (
	DefaultTimeFormat = "Mon, 02 Jan 2006 15:04:05.999 UTC"
)

var DefaultOptions = []Option{
	WithSkip(3),
	WithTimeFormat(DefaultTimeFormat),
	WithVerbosity(3),
	WithWriter(os.Stdout),
}

var DefaultLogger = NewFactory(DefaultOptions...).Logger

var Fatal = DefaultLogger(0)

var Error = DefaultLogger(1)

var Warning = DefaultLogger(2)

var Info = DefaultLogger(3)
