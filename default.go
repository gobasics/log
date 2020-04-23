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
	WithSkip(3),
	WithTimeFormat(defaultTimeFormat),
	WithVerbosity(defaultVerbosity),
	WithWriter(defaultWriter),
}

var V = NewFactory().Logger

var Err = V(0).Err

var Str = V(0).Str

var Strf = V(0).Strf
