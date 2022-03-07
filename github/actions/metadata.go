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
	Name        string             `json:"name,omitempty"`
	Author      string             `json:"author,omitempty"`
	Description string             `json:"description,omitempty"`
	Inputs      map[string]*Input  `json:"inputs,omitempty"`
	Outputs     map[string]*Output `json:"outputs,omitempty"`
	Runs        *Runs              `json:"runs,omitempty"`
}

type Input struct {
	Description        string `json:"input,omitempty"`
	Required           bool   `json:"required,omitempty"`
	Default            string `json:"default,omitempty"`
	DeprecationMessage string `json:"deprecationMessage,omitempty"`
}

type Output struct {
	Description string `json:"output,omitempty"`
}

type Runs struct {
	Plugin     string            `json:"plugin,omitempty"`
	Using      string            `json:"using,omitempty"`
	Main       string            `json:"main,omitempty"`
	Image      string            `json:"image,omitempty"`
	Entrypoint string            `json:"entrypoint,omitempty"`
	Args       []string          `json:"args,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
}
