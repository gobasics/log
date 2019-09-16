package log

type Level uint8

const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
)

var levels = []string{
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
}

func (level Level) String() string {
	var n = Level(len(levels))
	if level >= n {
		level = n - 1
	}

	return levels[level]
}
