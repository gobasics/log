package log

import "io"

type Option func(*Config)

func WithSkip(skip int) Option {
	return func(c *Config) {
		c.skip = skip
	}
}

func WithTimeFormat(tf string) Option {
	return func(c *Config) {
		c.timeFormat = tf
	}
}

func WithVerbosity(v uint8) Option {
	return func(c *Config) {
		c.verbosity = Level(v)
	}
}

func WithWriter(w io.Writer) Option {
	return func(c *Config) {
		c.writer = w
	}
}
