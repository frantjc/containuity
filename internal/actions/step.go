package actions

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/env"
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
	case "docker":
		if strings.HasPrefix("docker://", a.Runs.Image) {
			s.Image = a.Runs.Image
			s.Entrypoint = []string{a.Runs.Entrypoint}
			s.Cmd = a.Runs.Args
			s.Env = env.MapToArr(a.Runs.Env)
		} else {
			return nil, fmt.Errorf("action runs.using 'docker' only implemented for runs.image with prefix 'docker://'")
		}
	default:
		return nil, fmt.Errorf("action runs.using only implemented for 'node12', 'node16' and 'docker'")
	}

	return s, nil
}
