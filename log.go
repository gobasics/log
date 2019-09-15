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

type V interface {
	Err(error)
	Str(string)
}

type Logger interface {
	V(Level) V
}

type log struct {
	exit       func(exitCode int)
	fatal      func(level Level) bool
	fields     map[string][]interface{}
	levels     func(level Level) string
	verbose    func(level Level) bool
	timeFormat string
	writer     func(level Level) io.Writer
}

type logger struct {
	err func(error)
	str func(string)
}

func (l logger) Err(err error) { l.err(err) }

func (l logger) Str(str string) { l.str(str) }

func (l *log) V(level Level) V {
	return logger{
		err: func(err error) {
			if !l.verbose(level) {
				return
			}

			var message string

			if err != nil {
				message = err.Error()
			}

			l.log(level, message)
		},
		str: func(message string) {
			if !l.verbose(level) {
				return
			}
			l.log(level, message)
		},
	}
}

type logEntry struct {
	Fields    map[string][]interface{} `json:"fields,omitempty"`
	File      string                   `json:"file"`
	Level     string                   `json:"level"`
	Message   string                   `json:"message"`
	CreatedAt string                   `json:"createdAt"`
}

func (l *log) log(level Level, message string) {
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

type Option func(*log)

var filter = func(verbosity Level) func(level Level) bool {
	return func(level Level) bool {
		return level <= verbosity
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
			var verbosity = DEBUG
			var f = flag.NewFlagSet("gobasics/log", flag.ExitOnError)
			f.IntVar(&v, "v", int(verbosity), "Log verbosity level.")
			if err := f.Parse(os.Args); err == nil {
				verbosity = Level(v)
			}
			return filter(verbosity)(level)
		},
	),
}

func New(options ...Option) Logger {
	options = append(
		defaultOptions,
		options...,
	)
	var l = log{
		fields: make(map[string][]interface{}),
	}

	for _, o := range options {
		o(&l)
	}

	return &l
}
