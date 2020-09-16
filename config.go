package log

import (
	"io"
)

type Config struct {
	skip       int
	fields     map[string][]string
	timeFormat string
	verbosity  Level
	writer     io.Writer
}

func (c *Config) Apply(options ...Option) {
	for _, fn := range options {
		fn(c)
	}
}
