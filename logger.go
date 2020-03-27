package log

import (
	"fmt"
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
