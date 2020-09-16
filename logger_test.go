package log

import (
	"encoding/json"
	"errors"
	"strconv"
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

func (b *buffer) LastLog() *Log {
	var l Log
	line := b.LastLine()
	if line == nil {
		return nil
	}
	_ = json.Unmarshal(line, &l)
	return &l
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

func TestErr(t *testing.T) {
	for k, test := range []struct {
		err   error
		Level Level
		Log   string
	}{
		{nil, 0, ""},
		{errors.New("foo"), 0, "foo"},
		{errors.New("bar"), 0, "bar"},
	} {
		t.Run(strconv.Itoa(k), func(t *testing.T) {
			var w = buffer{}

			New(WithWriter(&w))(test.Level).Err(test.err)

			got := w.LastLog()

			if got != nil && test.Level != got.Level {
				t.Errorf("want %d, got %d", test.Level, got.Level)
			}

			if got != nil && test.Log != got.Log {
				t.Errorf("want %s, got %s", test.Log, got.Log)
			}
		})
	}
}

func TestStr(t *testing.T) {
	for k, v := range []struct {
		Level Level
		Log   string
	}{
		{0, "foo"},
		{1, "bar"},
	} {
		t.Run(strconv.Itoa(k), func(t *testing.T) {
			var w = buffer{}
			New(WithWriter(&w))(v.Level).Str(v.Log)

			got := w.LastLog()

			if got != nil && v.Level != got.Level {
				t.Errorf("want level=%d, got %d", v.Level, got.Level)
			}

			if got != nil && v.Log != got.Log {
				t.Errorf("want message=%s, got %s", v.Log, got.Log)
			}
		})
	}
}

func TestStrf(t *testing.T) {

	type logEntry struct {
		Level   Level  `json:"level"`
		Message string `json:"log"`
	}

	for k, test := range []struct {
		Level Level
		Log   string
	}{
		{0, "foo"},
		{1, "bar"},
	} {
		t.Run(strconv.Itoa(k), func(t *testing.T) {
			var w = buffer{}
			New(WithWriter(&w))(test.Level).Strf("%s", test.Log)

			got := w.LastLog()

			if got != nil && test.Level != got.Level {
				t.Errorf("want %d, got %d", test.Level, got.Level)
			}

			if got != nil && test.Log != got.Log {
				t.Errorf("want %s, got %s", test.Log, got.Log)
			}
		})
	}
}
