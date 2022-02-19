package actions

import (
	"io"

	"gopkg.in/yaml.v3"
)

func NewMetadataFromReader(r io.Reader) (*Metadata, error) {
	m := &Metadata{}
	d := yaml.NewDecoder(r)
	return m, d.Decode(m)
}

type Metadata struct {
	Name        string             `json:",omitempty"`
	Author      string             `json:",omitempty"`
	Description string             `json:",omitempty"`
	Inputs      map[string]*Input  `json:",omitempty"`
	Outputs     map[string]*Output `json:",omitempty"`
	Runs        *Runs              `json:",omitempty"`
}

type Input struct {
	Description        string `json:",omitempty"`
	Required           bool   `json:",omitempty"`
	Default            string `json:",omitempty"`
	DeprecationMessage string `json:",omitempty"`
}

type Output struct {
	Description string `json:",omitempty"`
}

type Runs struct {
	Plugin     string            `json:",omitempty"`
	Using      string            `json:",omitempty"`
	Main       string            `json:",omitempty"`
	Image      string            `json:",omitempty"`
	Entrypoint string            `json:",omitempty"`
	Args       []string          `json:",omitempty"`
	Env        map[string]string `json:",omitempty"`
}
