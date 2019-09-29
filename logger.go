package log

import "fmt"

type Logger interface {
	Err(error)
	Str(string)
	Strf(string, ...interface{})
}

type logger func(string)

func (l logger) Err(err error) {
	var s string
	if err != nil {
		s = err.Error()
	}
	l(s)
}

func (l logger) Str(s string) {
	l(s)
}

func (l logger) Strf(s string, args ...interface{}) {
	l(fmt.Sprintf(s, args...))
}
