package sequence

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"gopkg.in/yaml.v3"
)

const (
	imagePrefix = "docker://"
)

func NewStepFromReader(r io.Reader) (*Step, error) {
	s := &Step{}
	d := yaml.NewDecoder(r)
	return s, d.Decode(s)
}

func NewStepFromAction(a *actions.Action, path string) (*Step, error) {
	s := &Step{With: map[string]string{}}
	for inputName, input := range *a.Inputs {
		s.With[inputName] = input.Default
	}
	switch a.Runs.Using {
	case "node12":
		s.Image = "node:12"
		s.Entrypoint = []string{"node"}
		s.Cmd = []string{filepath.Join(path, a.Runs.Main)}
	case "node16":
		s.Image = "node:16"
		s.Entrypoint = []string{"node"}
		s.Cmd = []string{filepath.Join(path, a.Runs.Main)}
	case "docker":
		if strings.HasPrefix(a.Runs.Image, imagePrefix) {
			s.Image = strings.TrimPrefix(a.Runs.Image, imagePrefix)
			s.Entrypoint = []string{a.Runs.Entrypoint}
			s.Cmd = a.Runs.Args
			s.Env = a.Runs.Env
		} else {
			return nil, fmt.Errorf("action runs.using 'docker' only implemented for runs.image with prefix '%s'", imagePrefix)
		}
	default:
		return nil, fmt.Errorf("action runs.using only implemented for 'node12', 'node16' and 'docker'")
	}

	return s, nil
}

type Step struct {
	ID         string            `json:",omitempty"`
	Name       string            `json:",omitempty"`
	Image      string            `json:",omitempty"`
	Entrypoint []string          `json:",omitempty"`
	Cmd        []string          `json:",omitempty"`
	Privileged bool              `json:",omitempty"`
	Env        map[string]string `json:",omitempty"`

	Shell string            `json:",omitempty"`
	Run   string            `json:",omitempty"`
	Uses  string            `json:",omitempty"`
	With  map[string]string `json:",omitempty"`

	Get    string                 `json:",omitempty"`
	Put    string                 `json:",omitempty"`
	Params map[string]interface{} `json:",omitempty"`
}

func (s *Step) GetID() string {
	if s.ID != "" {
		return s.ID
	}
	return s.Name
}

func (s *Step) IsStdoutResponse() bool {
	return s.Uses != "" || s.Get != "" || s.Put != ""
}

func (s *Step) IsAction() bool {
	return s.Uses != "" || s.Run != ""
}

func (s *Step) Merge(step *Step) *Step {
	if s.ID == "" {
		s.ID = step.ID
	}
	if s.Name == "" {
		s.Name = step.Name
	}
	if s.Image == "" {
		s.Image = step.Image
	}
	if step.Privileged {
		s.Privileged = true
	}
	for key, value := range step.With {
		if s.With == nil {
			s.With = map[string]string{}
		}
		if s.With[key] == "" {
			s.With[key] = value
		}
	}

	return s
}
