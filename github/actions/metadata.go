package actions

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	RunsUsingDocker    = "docker"
	RunsUsingNode12    = "node12"
	RunsUsingNode16    = "node16"
	RunsUsingComposite = "composite"
)

func NewMetadataFromReader(r io.Reader) (*Metadata, error) {
	m := &Metadata{}
	d := yaml.NewDecoder(r)
	return m, d.Decode(m)
}

type Metadata struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Author      string `json:"author,omitempty" yaml:"author,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Inputs      map[string]*struct {
		Description        string      `json:"input,omitempty" yaml:"input,omitempty"`
		Required           bool        `json:"required,omitempty" yaml:"required,omitempty"`
		Default            interface{} `json:"default,omitempty" yaml:"default,omitempty"`
		DeprecationMessage string      `json:"deprecationMessage,omitempty" yaml:"deprecationMessage,omitempty"`
	} `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	Outputs map[string]*struct {
		Description string `json:"output,omitempty" yaml:"output,omitempty"`
	} `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Runs *struct {
		Plugin         string            `json:"plugin,omitempty" yaml:"plugin,omitempty"`
		Using          string            `json:"using,omitempty" yaml:"using,omitempty"`
		Pre            string            `json:"pre,omitempty" yaml:"pre,omitempty"`
		Main           string            `json:"main,omitempty" yaml:"main,omitempty"`
		Post           string            `json:"post,omitempty" yaml:"post,omitempty"`
		Image          string            `json:"image,omitempty" yaml:"image,omitempty"`
		PreEntrypoint  string            `json:"pre-entrypoint,omitempty" yaml:"pre-entrypoint,omitempty"`
		Entrypoint     string            `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
		PostEntrypoint string            `json:"post-entrypoint,omitempty" yaml:"post-entrypoint,omitempty"`
		Args           []string          `json:"args,omitempty" yaml:"args,omitempty"`
		Env            map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
		Steps          []*Step           `json:"steps,omitempty" yaml:"steps,omitempty"`
	} `json:"runs,omitempty" yaml:"runs,omitempty"`
}

func (m *Metadata) WithFromInputs() map[string]string {
	with := make(map[string]string, len(m.Inputs))
	for name, input := range m.Inputs {
		with[name] = fmt.Sprint(input.Default)
	}
	return with
}

func (m *Metadata) IsComposite() bool {
	return m.Runs.Using == RunsUsingComposite
}
