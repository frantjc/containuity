package planner

import (
	"context"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/github"
	"github.com/frantjc/sequence/pkg/container"
)

func defaultSpec() *container.Spec {
	return &container.Spec{
		Env: []string{"SEQUENCE=true"},
	}
}

func PlanStep(ctx context.Context, s *sequence.Step, opts ...PlanOpt) (*container.Spec, error) {
	popts := &planOpts{
		path: ".",
	}
	for _, opt := range opts {
		err := opt(popts)
		if err != nil {
			return nil, err
		}
	}

	spec := defaultSpec()
	spec.Privileged = s.Privileged

	if s.IsAction() {
		ghenv, err := github.NewEnv(popts.path)
		if err != nil {
			return nil, err
		}
		// TODO override step ID with ctx ID if exists in case this step is part of a job;
		// we want mounts to persist between steps in a job
		spec.Mounts = append(spec.Mounts, []container.Mount{
			{
				Source:      "/tmp/sqnc/workspace",
				Destination: ghenv.Workspace,
				Type:        container.MountTypeBind,
			},
			{
				Source:      "/tmp/sqnc/action",
				Destination: ghenv.ActionPath,
				Type:        container.MountTypeBind,
			},
			{
				Source:      "/tmp/sqnc/runner/temp",
				Destination: ghenv.RunnerTemp,
				Type:        container.MountTypeBind,
			},
			{
				Source:      "/tmp/sqnc/runner/toolcache",
				Destination: ghenv.RunnerToolCache,
				Type:        container.MountTypeBind,
			},
		}...)
		spec.Env = append(spec.Env, ghenv.Arr()...)
		spec.Cwd = ghenv.Workflow

		// s.Run doesn't necessarily need this image the way s.Uses does, but we may as well use it
		// since we own it and users will likely already have it stored locally
		spec.Image = sequence.Image()

		if s.Uses != "" {
			spec.Cmd = []string{"plugin", "uses", s.Uses, ghenv.ActionPath}
		} else if s.Run != "" {
			spec.Entrypoint = []string{"/bin/sh", "-c"}
			spec.Cmd = []string{s.Run}
		}
	} else {
		spec.Image = s.Image
		spec.Entrypoint = s.Entrypoint
		spec.Cmd = s.Cmd
	}

	return spec, nil
}
