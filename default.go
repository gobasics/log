package log

import (
	"os"
)

const (
	defaultTimeFormat = "Mon, 02 Jan 2006 15:04:05.999 UTC"
	defaultVerbosity  = 255
)

var defaultWriter = os.Stdout

var defaultOptions = []Option{
	WithTimeFormat(defaultTimeFormat),
	WithVerbosity(defaultVerbosity),
	WithWriter(defaultWriter),
}

var V = NewFactory().Logger
