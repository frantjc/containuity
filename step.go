package sequence

import (
	"io"

	"gopkg.in/yaml.v2"
)

func NewStepFromBytes(b []byte) (*Step, error) {
	s := &Step{}
	return s, yaml.Unmarshal(b, s)
}

func NewStepFromReader(r io.Reader) (*Step, error) {
	s := &Step{}
	d := yaml.NewDecoder(r)
	return s, d.Decode(s)
}

func NewStepFromString(s string) (*Step, error) {
	return NewStepFromBytes([]byte(s))
}

type Step struct {
	Name       string
	ID         string
	Image      string
	Entrypoint []string
	Cmd        []string
	Privileged bool

	Get string
	Put string

	Run  string
	Uses string
	With map[string]interface{}
}
