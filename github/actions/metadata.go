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
	Name        string             `json:"name,omitempty" yaml:"name,omitempty"`
	Author      string             `json:"author,omitempty" yaml:"author,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Inputs      map[string]*Input  `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	Outputs     map[string]*Output `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Runs        *Runs              `json:"runs,omitempty" yaml:"runs,omitempty"`
}

type Input struct {
	Description        string `json:"input,omitempty" yaml:"input,omitempty"`
	Required           bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Default            string `json:"default,omitempty" yaml:"default,omitempty"`
	DeprecationMessage string `json:"deprecationMessage,omitempty" yaml:"deprecationMessage,omitempty"`
}

type Output struct {
	Description string `json:"output,omitempty" yaml:"output,omitempty"`
}

type Runs struct {
	Plugin     string            `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	Using      string            `json:"using,omitempty" yaml:"using,omitempty"`
	Main       string            `json:"main,omitempty" yaml:"main,omitempty"`
	Image      string            `json:"image,omitempty" yaml:"image,omitempty"`
	Entrypoint string            `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Args       []string          `json:"args,omitempty" yaml:"args,omitempty"`
	Env        map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
}
