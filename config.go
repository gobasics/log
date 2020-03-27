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
	Fields       Fields
	TimeFormat   string
	VerboseLevel Level
	Writer       io.Writer
}

func (c Config) V(level Level) Logger {
	return c.Logger(level)
}

func (c Config) Logger(level Level) Logger {
	return func(message string) {
		if c.VerboseLevel >= level {
			c.log(level, message)
		}
	}
}

func (c Config) write(b []byte) {
	w := c.Writer
	if w != nil {
		_, _ = w.Write(b)
	}
}

func (config Config) log(level Level, message string) {
	type Log struct {
		Fields    Fields `json:"fields,omitempty"`
		File      string `json:"file"`
		Level     Level  `json:"level"`
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
	}
	var e = Log{
		Fields:    config.Fields,
		Level:     level,
		Message:   message,
		CreatedAt: time.Now().UTC().Format(config.TimeFormat),
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		e.File = fmt.Sprintf("%s:%d", file, line)
	}
	b, _ := json.Marshal(&e)
	b = append(b, '\n')
	config.write(b)
}
