package sequence

import (
	"fmt"
	"io"

	"github.com/google/go-containerregistry/pkg/name"
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
	Name       string   `yaml:"name"`
	IDF        string   `yaml:"id"`
	ImageF     string   `yaml:"image"`
	Entrypoint []string `yaml:"entrypoint"`
	Cmd        []string `yaml:"cmd"`
	Privileged bool     `yaml:"privileged"`

	Get string `yaml:"get"`
	Put string `yaml:"put"`

	Run  string                 `yaml:"run"`
	Uses string                 `yaml:"uses"`
	With map[string]interface{} `yaml:"with"`
}

var (
	_ Steppable = &Step{}
)

func (s *Step) Steps() ([]Step, error) {
	if s != nil {
		return []Step{*s}, nil
	}

	return nil, fmt.Errorf("nil step")
}

func (s *Step) ID() string {
	if s.IDF != "" {
		return s.IDF
	}

	return s.Name
}

func (s *Step) Validate() error {
	if _, err := s.Image(); err != nil {
		return err
	}

	if s.ID() == "" {
		return fmt.Errorf("one of id or name must be set")
	}

	return nil
}

func (s *Step) Image() (string, error) {
	var refs string
	if s.ImageF != "" {
		refs = s.ImageF
	}

	ref, err := name.ParseReference(refs)
	if err != nil {
		return "", fmt.Errorf("unable to parse image ref %s", refs)
	}

	return ref.Name(), nil
}
