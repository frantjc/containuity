package sequence

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// IsGitHubAction returns whether or not the step is a GitHub Action
func (s *Step) IsGitHubAction() bool {
	return s.Uses != ""
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
