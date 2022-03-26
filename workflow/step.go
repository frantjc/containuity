package workflow

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"gopkg.in/yaml.v3"
)

const (
	imagePrefix = "docker://"
	node12      = "docker.io/library/node:12"
	node16      = "docker.io/library/node:16"

	// ActionMetadataKey is the key in a Step's StepOut.Metadata
	// map that holds the json encoding of the action that
	// the step cloned
	ActionMetadataKey = "_sqnc_action"
)

// Step is a user's primary way of interacting with sequence;
// a Step defines a containerized command to be ran
// (or two if the Step is using a GitHub action from a GitHub repository:
//  first to clone the action and get its metadata, second to execute it)
type Step struct {
	ID    string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name  string            `json:"name,omitempty" yaml:"name,omitempty"`
	Env   map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Shell string            `json:"shell,omitempty" yaml:"shell,omitempty"`
	Run   string            `json:"run,omitempty" yaml:"run,omitempty"`
	Uses  string            `json:"uses,omitempty" yaml:"uses,omitempty"`
	With  map[string]string `json:"with,omitempty" yaml:"with,omitempty"`
	If    interface{}       `json:"if,omitempty" yaml:"if,omitempty"`

	Image      string   `json:"image,omitempty" yaml:"image,omitempty"`
	Entrypoint []string `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Cmd        []string `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Privileged bool     `json:"privileged,omitempty" yaml:"privileged,omitempty"`

	Get    string                 `json:"get,omitempty" yaml:"get,omitempty"`
	Put    string                 `json:"put,omitempty" yaml:"put,omitempty"`
	Params map[string]interface{} `json:"params,omitempty" yaml:"params,omitempty"`
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

func NewStepFromFile(name string) (*Step, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return NewStepFromReader(f)
}

// NewStepFromReader parses and returns a Step
// from the given reader e.g. a file
func NewStepFromReader(r io.Reader) (*Step, error) {
	s := &Step{}
	d := yaml.NewDecoder(r)
	return s, d.Decode(s)
}

// NewStepsFromMetadata returns a Step from a given GitHub action
// that is cloned at the given path
func NewStepsFromMetadata(a *actions.Metadata, path string) ([]Step, error) {
	var (
		steps = []Step{}
		with  = map[string]string{}
	)
	for inputName, input := range a.Inputs {
		with[inputName] = input.Default
	}
	switch a.Runs.Using {
	case "node12", "node16":
		image := node12
		if a.Runs.Using == "node16" {
			image = node16
		}
		if a.Runs.Pre != "" {
			steps = append(
				steps,
				Step{
					Image:      image,
					Entrypoint: []string{"node"},
					Cmd:        []string{filepath.Join(path, a.Runs.Pre)},
					With:       with,
					Env:        a.Runs.Env,
				},
			)
		}
		if a.Runs.Main != "" {
			steps = append(
				steps,
				Step{
					Image:      image,
					Entrypoint: []string{"node"},
					Cmd:        []string{filepath.Join(path, a.Runs.Main)},
					With:       with,
					Env:        a.Runs.Env,
				},
			)
		}
		if a.Runs.Post != "" {
			steps = append(
				steps,
				Step{
					Image:      image,
					Entrypoint: []string{"node"},
					Cmd:        []string{filepath.Join(path, a.Runs.Post)},
					With:       with,
					Env:        a.Runs.Env,
				},
			)
		}
	case "docker":
		if strings.HasPrefix(a.Runs.Image, imagePrefix) {
			image := strings.TrimPrefix(a.Runs.Image, imagePrefix)
			if entrypoint := a.Runs.PreEntrypoint; entrypoint != "" {
				steps = append(
					steps,
					Step{
						Image:      image,
						Entrypoint: []string{entrypoint},
						With:       with,
						Env:        a.Runs.Env,
					},
				)
			}
			if entrypoint := a.Runs.Entrypoint; entrypoint != "" {
				steps = append(
					steps,
					Step{
						Image:      image,
						Entrypoint: []string{entrypoint},
						Cmd:        a.Runs.Args,
						With:       with,
						Env:        a.Runs.Env,
					},
				)
			}
			if entrypoint := a.Runs.PostEntrypoint; entrypoint != "" {
				steps = append(
					steps,
					Step{
						Image:      image,
						Entrypoint: []string{entrypoint},
						With:       with,
						Env:        a.Runs.Env,
					},
				)
			}
		} else {
			return nil, fmt.Errorf("action runs.using 'docker' only implemented for runs.image with prefix '%s', got '%s'", imagePrefix, a.Runs.Image)
		}
	default:
		return nil, fmt.Errorf("action runs.using only implemented for 'node12', 'node16' and 'docker', got '%s'", a.Runs.Using)
	}

	return steps, nil
}

// StepOut is the optional parsable output of a Step
// on its stdout
// e.g. if a Step is a Concourse Resource or a sqncshim
// that is returning metadata about the action that it cloned
type StepOut struct {
	Metadata map[string]string      `json:"metadata,omitempty"`
	Version  map[string]interface{} `json:"version,omitempty"`
}
