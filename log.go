package log

import "encoding/json"

type Log struct {
	Level  Level               `json:"level"`
	At     string              `json:"at"`
	Log    string              `json:"log"`
	File   string              `json:"file,omitempty"`
	Fields map[string][]string `json:"fields,omitempty"`
}

func (l Log) Bytes() []byte {
	b, _ := json.Marshal(&l)
	b = append(b, 10)
	return b
}
