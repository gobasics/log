package log

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"time"
)

type Fields map[string][]string

type Level uint8

type Factory struct {
	skip       int
	fields     Fields
	timeFormat string
	verbosity  Level
	writer     io.Writer
}

func (f Factory) Logger(level Level) Logger {
	return func(message string) {
		if f.verbosity >= level {
			f.log(level, message)
		}
	}
}

func (f Factory) write(b []byte) {
	w := f.writer
	_, _ = w.Write(b)
}

func (f Factory) log(level Level, message string) {
	type Log struct {
		Level  Level  `json:"level"`
		At     string `json:"at"`
		Log    string `json:"log"`
		File   string `json:"file"`
		Fields Fields `json:"fields,omitempty"`
	}
	var e = Log{
		Fields: f.fields,
		Level:  level,
		Log:    message,
		At:     time.Now().UTC().Format(f.timeFormat),
	}
	if _, file, line, ok := runtime.Caller(f.skip); ok {
		e.File = fmt.Sprintf("%s:%d", file, line)
	}
	b, _ := json.Marshal(&e)
	b = append(b, '\n')
	f.write(b)
}

func NewFactory(options ...Option) Factory {
	var f Factory
	for _, fn := range DefaultOptions {
		fn(&f)
	}
	for _, fn := range options {
		fn(&f)
	}
	return f
}
