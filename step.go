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
	Name       string   `json:",omitempty"`
	IDF        string   `json:"id,omitempty"`
	Image      string   `json:",omitempty"`
	Entrypoint []string `json:",omitempty"`
	Cmd        []string `json:",omitempty"`
	Privileged bool     `json:",omitempty"`
	Env        []string `json:",omitempty"`

	Run  string                 `json:",omitempty"`
	Uses string                 `json:",omitempty"`
	With map[string]interface{} `json:",omitempty"`

	Get string `json:",omitempty"`
	Put string `json:",omitempty"`
}

func (s *Step) ID() string {
	if s.IDF != "" {
		return s.IDF
	}
	return s.Name
}

func (s *Step) IsStdoutParsable() bool {
	return s.Uses != "" || s.Get != "" || s.Put != ""
}

func (s *Step) IsAction() bool {
	return s.Uses != "" || s.Run != ""
}
