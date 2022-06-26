package sequence

import (
	"io"
	"os"

	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v3"
)

// GetID returns the functional ID of the step
func (s *Step) GetID() string {
	if s.Id != "" {
		return s.Id
	} else if s.Name != "" {
		return s.Name
	}
	return s.Uses
}

// IsStdoutResponse returns whether or not this step is expected to
// respond with a StepResponse on stdout or not
func (s *Step) IsStdoutResponse() bool {
	return s.IsGitHubAction() || s.IsConcourseResource()
}

// IsAction returns whether or not this step is a GitHub Action
func (s *Step) IsGitHubAction() bool {
	return s.Uses != ""
}

// IsResource returns whether or not this step is a Concourse Resoucre
func (s *Step) IsConcourseResource() bool {
	return s.Get != "" || s.Put != ""
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
	if s.Id == "" {
		s.Id = step.Id
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
	if step.Id != "" {
		s.Id = step.Id
	}
	if step.Name != "" {
		s.Name = step.Name
	}
	if step.Image != "" {
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
		s.Params = map[string]*anypb.Any{}
	} else if s.Uses != "" {
		s.Get = ""
		s.Put = ""
		s.Params = map[string]*anypb.Any{}
	}

	return s
}

// NewStepFromFile parses and returns a Step
// from the given file name
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
