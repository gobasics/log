package log

type Fields map[string][]string

func (f Fields) Add(k, v string) {
	f[k] = append(f[k], v)
}
