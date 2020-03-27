package log

import (
	"encoding/json"
	"errors"
	"fmt"
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

func equal(f func(string, ...interface{})) func(string) func(interface{}, interface{}) {
	return func(message string) func(interface{}, interface{}) {
		return func(got, want interface{}) {
			if got != want {
				f(message, got, want)
			}
		}
	}
}

func TestErr(t *testing.T) {
	const message = "TestErr"
	var err = errors.New(message)

	type logEntry struct {
		Level   Level  `json:"level"`
		Message string `json:"message"`
	}

	for _, test := range []struct {
		name  string
		err   error
		Level Level

		Entry logEntry
	}{
		{"a", nil, 0, logEntry{Level: 0, Message: ""}},
		{"b", err, 0, logEntry{Level: 0, Message: message}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var w = buffer{}

			var p = Config{
				TimeFormat: defaultTimeFormat,
				Verbosity:  defaultWriteLevel,
				Writer:     &w,
				Fields:     make(Fields),
			}

			p.Fields.Add("foo", "bar")

			var gotEntry logEntry

			p.V(test.Level).Err(test.err)

			if line := w.LastLine(); line != nil {
				err = json.Unmarshal(line, &gotEntry)
				equal(t.Errorf)("err -> got=%#v, want=%#v")(err, nil)
			}

			var eq = equal(t.Errorf)

			eq("level -> got=%d, want=%d")(gotEntry.Level, test.Entry.Level)

			eq("message -> got=%s, want=%s")(gotEntry.Message, test.Entry.Message)
		})
	}
}

func TestStr(t *testing.T) {
	type logEntry struct {
		Level   Level  `json:"level"`
		Message string `json:"message"`
	}

	for _, test := range []struct {
		name  string
		Level Level
	}{
		{"a", 0},
		{"b", 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			const wantMessage = "TestStr"

			var w = buffer{}
			var l = Config{
				TimeFormat: defaultTimeFormat,
				Verbosity:  defaultWriteLevel,
				Writer:     &w,
			}
			l.V(test.Level).Str(wantMessage)

			var got logEntry
			var line = w.LastLine()
			if line != nil {
				var err = json.Unmarshal(line, &got)
				equal(t.Errorf)("err -> got=%#v, want=%#v; %w")(err, nil)
			}

			var eq = equal(t.Errorf)

			eq("level -> got=%d, want=%d")(got.Level, test.Level)

			eq("message -> got=%s, want=%s")(got.Message, wantMessage)
		})
	}
}

func TestStrf(t *testing.T) {

	type logEntry struct {
		Level   Level  `json:"level"`
		Message string `json:"message"`
	}

	for _, test := range []struct {
		name  string
		Level Level
	}{
		{"a", 0},
		{"b", 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			const wrapMessage = "TestStr; %s"

			var cause = "cause"

			var wantMessage = fmt.Sprintf(wrapMessage, cause)

			var w = buffer{}
			var l = Config{
				TimeFormat: defaultTimeFormat,
				Verbosity:  defaultWriteLevel,
				Writer:     &w,
			}
			l.V(test.Level).Strf(wrapMessage, cause)

			var got logEntry
			var line = w.LastLine()
			if line != nil {
				var err = json.Unmarshal(line, &got)
				equal(t.Errorf)("err -> got=%#v, want=%#v")(err, nil)
			}

			var eq = equal(t.Errorf)

			eq("level -> got=%d, want=%d")(got.Level, test.Level)

			eq("message -> got=%s, want=%s")(got.Message, wantMessage)
		})
	}
}
