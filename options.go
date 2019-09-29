package log

import (
	"flag"
	"io"
	"os"
)

type Option func(*provider)

func WithExitFunc(f func(exitCode int)) Option {
	return func(p *provider) {
		p.exit = f
	}
}

func WithField(key string, value interface{}) Option {
	return func(p *provider) {
		values := p.fields[key]
		p.fields[key] = append(values, value)
	}
}

func WithTimeFormat(tf string) Option {
	return func(p *provider) {
		p.timeFormat = tf
	}
}

func WithDefaultWriters() Option {
	var (
		Fatal   io.Writer = os.Stderr
		Error             = os.Stderr
		Warning           = os.Stderr
		Info              = os.Stdout
		Debug             = os.Stdout
	)
	return WithWriters(Fatal, Error, Warning, Info, Debug)
}

func WithWriters(writers ...io.Writer) Option {
	return func(p *provider) {
		p.writer = func(level Level) io.Writer {
			if level < Level(len(writers)) {
				return writers[level]
			}
			return writers[len(writers)-1]
		}
	}
}

func WithVerbosity(max Level) Option {
	return func(p *provider) {
		p.verbose = func(level Level) bool {
			return level <= max
		}
	}
}

func WithVerbosityFlag() Option {
	var verbosity = func() Level {
		var v int
		var max = DEBUG
		var f = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		f.IntVar(&v, "v", int(max), "Log verbosity level.")
		if err := f.Parse(os.Args); err == nil {
			max = Level(v)
		}
		return max
	}
	return WithVerbosity(verbosity())
}
