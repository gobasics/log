package log

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

type Logger interface {
	Err(error)
	Str(string)
	Strf(string, ...interface{})
}

type logger struct {
	exit       func(exitCode int)
	fatal      func(level Level) bool
	fields     map[string][]interface{}
	levels     func(level Level) string
	verbose    func(level Level) bool
	timeFormat string
	writer     func(level Level) io.Writer
}

type verbosity func(message string)

func (v verbosity) Err(err error) {
	var message string
	if err != nil {
		message = err.Error()
	}
	v(message)
}

func (v verbosity) Str(message string) {
	v(message)
}

func (l *logger) V(level Level) interface {
	Err(error)
	Str(string)
} {
	var v verbosity = func(message string) {
		if l.verbose(level) {
			l.log(level, message)
		}
	}
	return v
}

type logEntry struct {
	Fields    map[string][]interface{} `json:"fields,omitempty"`
	File      string                   `json:"file"`
	Level     string                   `json:"level"`
	Message   string                   `json:"message"`
	CreatedAt string                   `json:"createdAt"`
}

func (l *logger) log(level Level, message string) {
	var e = logEntry{
		CreatedAt: time.Now().Format(l.timeFormat),
		Fields:    l.fields,
		Level:     l.levels(level),
		Message:   message,
	}

	if _, file, line, ok := runtime.Caller(2); ok {
		e.File = fmt.Sprintf("%s:%d", file, line)
	}

	b, _ := json.Marshal(&e)

	l.writer(level).Write(b)

	if l.fatal(level) {
		l.exit(1)
	}
}

type Option func(*logger)

var filter = func(max Level) func(level Level) bool {
	return func(level Level) bool {
		return level <= max
	}
}

var selectWriter = func(writers ...io.Writer) func(level Level) io.Writer {
	return func(level Level) io.Writer {
		var w io.Writer = os.Stdout
		if level < Level(len(writers)) {
			w = writers[level]
		}
		return w
	}
}

var defaultOptions = []Option{
	WithExitFunc(os.Exit),
	WithFatalFunc(func(level Level) bool { return level == FATAL }),
	WithField("foo", "bar"),
	WithLevelsFunc(
		func(level Level) string {
			return level.String()
		},
	),
	WithTimeFormat("Mon, 02 Jan 2006 15:04:05.999 UTC"),
	WithWriterFunc(
		selectWriter(os.Stderr, os.Stderr, os.Stderr, os.Stdout),
	),
	WithVerbosityFunc(
		func(level Level) bool {
			var v int
			var max = DEBUG
			var f = flag.NewFlagSet("gobasics/log", flag.ExitOnError)
			f.IntVar(&v, "v", int(max), "Log verbosity level.")
			if err := f.Parse(os.Args); err == nil {
				max = Level(v)
			}
			return filter(max)(level)
		},
	),
}

func New(options ...Option) func(Level) Logger {
	options = append(
		defaultOptions,
		options...,
	)
	var l = logger{
		fields: make(map[string][]interface{}),
	}

	for _, o := range options {
		o(&l)
	}

	return &l
}
