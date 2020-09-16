package log

import (
	"fmt"
	"runtime"
	"time"
)

type Logger func(string)

func (l Logger) Err(err error) {
	if err != nil {
		l(err.Error())
	}
}

func (l Logger) Str(s string) {
	l(s)
}

func (l Logger) Strf(s string, args ...interface{}) {
	l(fmt.Sprintf(s, args...))
}

func New(options ...Option) func(Level) Logger {
	var c Config

	(&c).Apply(options...)

	return func(level Level) Logger {
		return func(message string) {
			if c.verbosity < level {
				return
			}

			var e = Log{
				Fields: c.fields,
				Level:  level,
				Log:    message,
				At:     time.Now().UTC().Format(c.timeFormat),
			}

			if _, file, line, ok := runtime.Caller(c.skip); ok {
				e.File = fmt.Sprintf("%s:%d", file, line)
			}

			_, _ = c.writer.Write(e.Bytes())
		}
	}
}
