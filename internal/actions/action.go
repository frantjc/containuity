package actions

import (
	"io"

	"gopkg.in/yaml.v2"
)

func NewActionFromBytes(b []byte) (*Action, error) {
	a := &Action{}
	return a, yaml.Unmarshal(b, a)
}

func NewActionFromReader(r io.Reader) (*Action, error) {
	s := &Action{}
	d := yaml.NewDecoder(r)
	return s, d.Decode(s)
}

func NewActionFromString(s string) (*Action, error) {
	return NewActionFromBytes([]byte(s))
}

type Action struct {
	Name        string
	Author      string
	Description string
	Inputs      *Inputs
	Outputs     *Outputs
	Runs        *Runs
}

type Inputs map[string]struct {
	Description        string
	Required           bool
	Default            string
	DeprecationMessage string
}

type Outputs map[string]struct {
	Description string
}

type Runs struct {
	Plugin     string
	Using      string
	Main       string
	Image      string
	Entrypoint string
	Args       []string
	Env        map[string]string
}
