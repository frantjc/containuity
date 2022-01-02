package actions

import (
	"fmt"
	"path/filepath"

	"github.com/frantjc/sequence"
)

func (a *Action) Step(path string) (*sequence.Step, error) {
	return ToStep(a, path)
}

func ToStep(a *Action, path string) (*sequence.Step, error) {
	s := &sequence.Step{}
	switch a.Runs.Using {
	case "node12":
		s.Image = "node:12"
		s.Entrypoint = []string{"node"}
		s.Cmd = []string{filepath.Join(path, a.Runs.Main)}
	case "node16":
		s.Image = "node:16"
		s.Entrypoint = []string{"node"}
		s.Cmd = []string{filepath.Join(path, a.Runs.Main)}
	default:
		return nil, fmt.Errorf("action runs.using only implemented for 'node12' and 'node16'")
	}

	return s, nil
}
