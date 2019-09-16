package log

import (
	"io"
)

func WithExitFunc(f func(exitCode int)) Option {
	return func(l *logger) {
		l.exit = f
	}
}

func WithFatalFunc(f func(level Level) bool) Option {
	return func(l *logger) {
		l.fatal = f
	}
}

func WithField(key string, value interface{}) Option {
	return func(l *logger) {
		values := l.fields[key]
		l.fields[key] = append(values, value)
	}
}

func WithLevelsFunc(f func(level Level) (levelName string)) Option {
	return func(l *logger) {
		l.levels = f
	}
}

func WithTimeFormat(tf string) Option {
	return func(l *logger) {
		l.timeFormat = tf
	}
}

func WithWriterFunc(f func(level Level) io.Writer) Option {
	return func(l *logger) {
		l.writer = f
	}
}

func WithVerbosityFunc(f func(level Level) bool) Option {
	return func(l *logger) {
		l.verbose = f
	}
}
