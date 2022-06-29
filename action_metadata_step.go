package sequence

import (
	errors "errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/github/actions"
)

var (
	ErrCompositeActionMetadata    = errors.New("action metadata is composite")
	ErrNotCompositeActionMetadata = errors.New("action metadata is not composite")
)

func NewStepsFromCompositeActionMetadata(a *actions.Metadata, path string) ([]*Step, error) {
	if !a.IsComposite() {
		return nil, ErrNotCompositeActionMetadata
	}

	steps := []*Step{}
	for _, step := range a.Runs.Steps {
		steps = append(steps, &Step{
			Env:   step.Env,
			Id:    step.ID,
			If:    step.If,
			Name:  step.Name,
			Run:   step.Run,
			Shell: step.Shell,
			Uses:  step.Uses,
			With:  step.With,
			// TODO WorkingDirectory
		})
	}

	return steps, nil
}

func NewStepsFromNonCompositeMetadata(a *actions.Metadata, path string, parentStep *Step) (*Step, *Step, *Step, error) {
	if a.IsComposite() {
		return nil, nil, nil, ErrCompositeActionMetadata
	}

	var (
		preStep        *Step
		mainStep       *Step
		postStep       *Step
		image          = ImageNode12.GetRef()
		preEntrypoint  []string
		preCmd         []string
		mainEntrypoint []string
		mainCmd        []string
		postEntrypoint []string
		postCmd        []string
		with           = a.WithFromInputs()
	)
	switch a.Runs.Using {
	case actions.RunsUsingNode12, actions.RunsUsingNode16:
		if a.Runs.Using == actions.RunsUsingNode16 {
			image = ImageNode16.GetRef()
		}

		if a.Runs.Pre != "" {
			preEntrypoint = []string{"node"}
			preCmd = []string{filepath.Join(path, a.Runs.Pre)}
		}

		// every non-composite action has a main
		mainEntrypoint = []string{"node"}
		mainCmd = []string{filepath.Join(path, a.Runs.Main)}

		if a.Runs.Post != "" {
			postEntrypoint = []string{"node"}
			postCmd = []string{filepath.Join(path, a.Runs.Post)}
		}
	case actions.RunsUsingDocker:
		if strings.HasPrefix(a.Runs.Image, actions.RunsUsingDockerImagePrefix) {
			image = strings.TrimPrefix(a.Runs.Image, actions.RunsUsingDockerImagePrefix)
		} else {
			return nil, nil, nil, fmt.Errorf("action runs.using '%s' only implemented for runs.image with prefix '%s', got '%s'", actions.RunsUsingDocker, actions.RunsUsingDockerImagePrefix, a.Runs.Image)
		}

		if entrypoint := a.Runs.PreEntrypoint; entrypoint != "" {
			preEntrypoint = []string{entrypoint}
		}

		if entrypoint := a.Runs.Entrypoint; entrypoint != "" {
			mainEntrypoint = []string{entrypoint}
		}

		if entrypoint := a.Runs.PostEntrypoint; entrypoint != "" {
			postEntrypoint = []string{entrypoint}
		}
	default:
		return nil, nil, nil, fmt.Errorf("action runs.using only implemented for '%s', '%s' and '%s', got '%s'", actions.RunsUsingDocker, actions.RunsUsingNode12, actions.RunsUsingNode16, a.Runs.Using)
	}

	if len(preEntrypoint) > 0 {
		preStep = &Step{
			Id:         js.Ternary(parentStep.Id != "", fmt.Sprintf("Pre %s", parentStep.Id), ""),
			Name:       js.Ternary(parentStep.Name != "", fmt.Sprintf("Pre %s", parentStep.Name), ""),
			Image:      image,
			Entrypoint: preEntrypoint,
			Cmd:        preCmd,
			With:       with,
			Env:        a.Runs.Env,
		}
	}

	if len(mainEntrypoint) > 0 {
		mainStep = &Step{
			Id:         parentStep.Id,
			Name:       parentStep.Name,
			Image:      image,
			Entrypoint: mainEntrypoint,
			Cmd:        mainCmd,
			With:       with,
			Env:        a.Runs.Env,
		}
	}

	if len(postEntrypoint) > 0 {
		postStep = &Step{
			Id:         js.Ternary(parentStep.Id != "", fmt.Sprintf("Post %s", parentStep.Id), ""),
			Name:       js.Ternary(parentStep.Name != "", fmt.Sprintf("Post %s", parentStep.Name), ""),
			Image:      image,
			Entrypoint: postEntrypoint,
			Cmd:        postCmd,
			With:       with,
			Env:        a.Runs.Env,
		}
	}

	for _, childStep := range []*Step{
		preStep, mainStep, postStep,
	} {
		if childStep != nil {
			for k, v := range parentStep.With {
				childStep.With[k] = v
			}

			for k, v := range parentStep.Env {
				childStep.Env[k] = v
			}
		}
	}

	return preStep, mainStep, postStep, nil
}
