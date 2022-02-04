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

// GetID returns the functional ID of the step
func (s *Step) GetID() string {
	if s.ID != "" {
		return s.ID
	} else if s.Name != "" {
		return s.Name
	}
	return s.Uses
}

// IsStdoutResponse returns whether or not this step is expected to
// respond with a StepResponse on stdout or not
func (s *Step) IsStdoutResponse() bool {
	return s.Uses != "" || s.Get != "" || s.Put != ""
}

// IsStdoutResponse returns whether or not this step is a GitHub Action
// or not
func (s *Step) IsAction() bool {
	return s.Uses != "" || s.Run != ""
}

// Merge sets all of this step's undefined fields with
// the given step's fields
func (s *Step) Merge(step *Step) *Step {
	if s.Entrypoint == nil || len(s.Entrypoint) == 0 {
		s.Entrypoint = step.Entrypoint
	}
	if s.Cmd == nil || len(s.Cmd) == 0 {
		s.Cmd = step.Cmd
	}
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

// MergeOverride overrides all of this step's fields with
// the given step's fields if they are defined
func (s *Step) MergeOverride(step *Step) *Step {
	if step.ID == "" {
		s.ID = step.ID
	}
	if step.Name == "" {
		s.Name = step.Name
	}
	if step.Image == "" {
		s.Image = step.Image
	}
	if step.Privileged {
		s.Privileged = true
	}
	for key, value := range step.With {
		if s.With == nil {
			s.With = map[string]string{}
		}
		s.With[key] = value
	}

	return s
}

// Canonical returns the Step's "canonical" form, e.g. the step
//
// image: alpine
// uses: actions/checkout@v2
//
// is an oxymoron, so this function would make it into
//
// image: alpine
func (s *Step) Canonical() *Step {
	if s.Image != "" {
		s.Uses = ""
		s.Get = ""
		s.Put = ""
		s.Params = map[string]interface{}{}
	} else if s.Uses != "" {
		s.Get = ""
		s.Put = ""
		s.Params = map[string]interface{}{}
	}

	return s
}
