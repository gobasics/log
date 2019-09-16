package log

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"testing"
)

type buffer struct {
	data  [][]byte
	mutex sync.Mutex
}

func (b *buffer) Write(v []byte) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.data = append(b.data, v)
	return len(v), nil
}

func (b *buffer) LastLine() []byte {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	n := len(b.data)

	if n > 0 {
		return b.data[n-1]
	}

	return nil
}

func equalf(f func(string, ...interface{})) func(string) func(interface{}, interface{}) {
	return func(message string) func(interface{}, interface{}) {
		return func(got, want interface{}) {
			if got != want {
				f(message, got, want)
			}
		}
	}
}

func TestE(t *testing.T) {
	const message = "TestE"

	var err = errors.New(message)

	for _, test := range []struct {
		name   string
		err    error
		Exited bool
		Level  Level

		Entry logEntry
	}{
		{"a", nil, true, FATAL, logEntry{Level: FATAL.String(), Message: ""}},
		{"b", err, true, FATAL, logEntry{Level: FATAL.String(), Message: message}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var l = New()

			var gotExited bool
			var gotEntry logEntry

			var w = &buffer{}
			l.(*logger).writer = func(Level) io.Writer { return w }
			l.(*logger).exit = func(int) { gotExited = true }

			l.V(test.Level).Err(test.err)

			var line = w.LastLine()
			if line != nil {
				err = json.Unmarshal(line, &gotEntry)
				equalf(t.Fatalf)("err -> got=%#v, want=%#v")(err, nil)
			}

			var eq = equalf(t.Errorf)

			eq("exited -> got=%t, want=%t")(gotExited, test.Exited)

			eq("level -> got=%s, want=%s")(gotEntry.Level, test.Entry.Level)

			eq("message -> got=%s, want=%s")(gotEntry.Message, test.Entry.Message)
		})
	}
}

func TestNew(t *testing.T) {
	var exits []int

	for _, test := range []struct {
		name  string
		level Level
	}{
		{"a", FATAL},
		{"b", ERROR},
		{"c", WARN},
		{"d", INFO},
		{"e", DEBUG},
		{"f", 5},
		{"g", 6},
		{"h", 7},
		{"i", 8},
		{"j", 9},
	} {
		t.Run(
			test.name,
			func(t *testing.T) {
				var w = &buffer{}

				var options = []Option{
					WithExitFunc(func(code int) { exits = append(exits, code) }),
					WithWriterFunc(selectWriter(w, w, w, w, w, w, w, w, w, w)),
					WithVerbosityFunc(filter(5)),
				}

				var l = New(options...)
				l.V(test.level).Err(errors.New(test.name))
				l.V(test.level).Str(test.name)
			},
		)
	}
}

func TestS(t *testing.T) {
	const message = "TestS"

	var err = errors.New(message)

	for _, test := range []struct {
		name   string
		Exited bool
		Level  Level

		Log logEntry
	}{
		{"a", true, FATAL, logEntry{Level: FATAL.String(), Message: message}},
		{"b", false, ERROR, logEntry{Level: ERROR.String(), Message: message}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var w = buffer{}
			var gotExited bool
			var l = New(
				WithExitFunc(func(int) { gotExited = true }),
				WithWriterFunc(func(Level) io.Writer { return &w }),
			)

			var gotLog logEntry

			l.V(test.Level).Str(test.Log.Message)

			var line = w.LastLine()
			if line != nil {
				err = json.Unmarshal(line, &gotLog)
				equalf(t.Fatalf)("err -> got=%#v, want=%#v")(err, nil)
			}

			var eq = equalf(t.Errorf)

			eq("exited -> got=%t, want=%t")(gotExited, test.Exited)

			eq("level -> got=%s, want=%s")(gotLog.Level, test.Log.Level)

			eq("message -> got=%s, want=%s")(gotLog.Message, test.Log.Message)
		})
	}
}
