package log

import (
	"os"
)

const defaultTimeFormat = "Mon, 02 Jan 2006 15:04:05.999 UTC"

const defaultWriteLevel = 255

var defaultWriter = os.Stdout

var defaultConfig = Config{
	TimeFormat: defaultTimeFormat,
	Verbosity:  defaultWriteLevel,
	Writer:     defaultWriter,
}

var V = defaultConfig.Logger
