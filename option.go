package log

import "io"

type Option func(*Factory)

func WithVerbosity(v Level) Option {
	return func(c *Factory) {
		c.verbosity = v
	}
}
func WithWriter(w io.Writer) Option {
	return func(c *Factory) {
		c.writer = w
	}
}
func WithTimeFormat(tf string) Option {
	return func(c *Factory) {
		c.timeFormat = tf
	}
}
