package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFile(fn string) (b []byte, err error) {
	b, err = os.ReadFile(fn)
	return
}

func LoadYAML(b []byte, v interface{}) (err error) {
	return yaml.Unmarshal(b, v)
}
