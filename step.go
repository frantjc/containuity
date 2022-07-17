package sequence

import (
	"io"
	"os"

	"github.com/frantjc/go-js"
	"gopkg.in/yaml.v3"
)

// GetID returns the step's effective ID.
func (s *Step) GetID() string {
	return js.Coalesce(s.Id, s.Name)
}

// IsGitHubAction returns whether or not the step is a GitHub Action.
func (s *Step) IsGitHubAction() bool {
	return s.Uses != ""
}

// NewStepFromFile parses and returns a Step
// from the given file name.
func NewStepFromFile(name string) (*Step, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return NewStepFromReader(f)
}

// NewStepFromReader parses and returns a Step
// from the given reader e.g. a file.
func NewStepFromReader(r io.Reader) (*Step, error) {
	s := &Step{}
	d := yaml.NewDecoder(r)
	return s, d.Decode(s)
}
