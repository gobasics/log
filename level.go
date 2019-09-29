package log

type Level uint8

const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
	DEBUG

	WARNING = WARN
)

func (l Level) String() string {
	var levels = []string{
		"FATAL",
		"ERROR",
		"WARN",
		"INFO",
		"DEBUG",
	}

	var n = Level(len(levels))
	if l >= n {
		return levels[n-1]
	}

	return levels[l]
}
