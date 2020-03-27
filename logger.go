package log

import (
	"fmt"
)

type Logger func(string)

func (l Logger) Err(err error) {
	var s string
	if err != nil {
		s = err.Error()
	}
	l(s)
}

func (l Logger) Str(s string) {
	l(s)
}

func (l Logger) Strf(s string, args ...interface{}) {
	l(fmt.Sprintf(s, args...))
}
