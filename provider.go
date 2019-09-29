package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

type provider struct {
	exit       func(exitCode int)
	fields     map[string][]interface{}
	verbose    func(level Level) bool
	timeFormat string
	writer     func(level Level) io.Writer
}

func (p *provider) V(level Level) Logger {
	var l logger = func(message string) {
		if p.verbose(level) {
			p.log(level, message)
		}
	}
	return l
}

func (p *provider) log(level Level, message string) {
	type logEntry struct {
		Fields    map[string][]interface{} `json:"fields,omitempty"`
		File      string                   `json:"file"`
		Level     string                   `json:"level"`
		Message   string                   `json:"message"`
		CreatedAt string                   `json:"createdAt"`
	}

	var e = logEntry{
		CreatedAt: time.Now().Format(p.timeFormat),
		Fields:    p.fields,
		Level:     level.String(),
		Message:   message,
	}

	if _, file, line, ok := runtime.Caller(1); ok {
		e.File = fmt.Sprintf("%s:%d", file, line)
	}

	b, _ := json.Marshal(&e)

	_, _ = p.writer(level).Write(b)

	if FATAL == level {
		p.exit(1)
	}
}

type Provider interface{ V(Level) Logger }

func New(options ...Option) Provider {
	var p = provider{
		fields: make(map[string][]interface{}),
	}

	const tf = "Mon, 02 Jan 2006 15:04:05.999 UTC"

	for _, o := range []Option{
		WithExitFunc(os.Exit),
		WithField("start_time", time.Now().UTC().Format(tf)),
		WithTimeFormat(tf),
		WithDefaultWriters(),
		WithVerbosityFlag(),
	} {
		o(&p)
	}

	for _, o := range options {
		o(&p)
	}

	return &p
}
