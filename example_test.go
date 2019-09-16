package log_test

import (
	"errors"

	"github.com/gobasics/log"
)

func ExampleNew() {
	var l = log.New()

	l.V(0).Str("foo")

	l.V(0).Err(errors.New("bar"))
}
