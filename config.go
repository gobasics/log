package log

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"time"
)

type Level uint8

type Config struct {
	Fields     Fields
	TimeFormat string
	Verbosity  Level
	Writer     io.Writer
}

func (c Config) V(level Level) Logger {
	return c.Logger(level)
}

func (c Config) Logger(level Level) Logger {
	return func(message string) {
		if c.Verbosity >= level {
			c.log(level, message)
		}
	}
}

func (c Config) write(b []byte) {
	w := c.Writer
	if w == nil {
		w = defaultWriter
	}
	_, _ = w.Write(b)
}

func (config Config) log(level Level, message string) {
	type Log struct {
		Level  Level  `json:"level"`
		At     string `json:"at"`
		Log    string `json:"log"`
		Fields Fields `json:"fields,omitempty"`
		File   string `json:"file"`
	}
	var e = Log{
		Fields: config.Fields,
		Level:  level,
		Log:    message,
		At:     config.timeNow(),
	}
	if _, file, line, ok := runtime.Caller(3); ok {
		e.File = fmt.Sprintf("%s:%d", file, line)
	}
	b, _ := json.Marshal(&e)
	b = append(b, '\n')
	config.write(b)
}

func (c Config) timeNow() string {
	tf := c.TimeFormat
	if tf == "" {
		tf = defaultTimeFormat
	}
	return time.Now().UTC().Format(tf)
}
