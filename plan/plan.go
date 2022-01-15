package plan

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/env"
	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func defaultSpec() *runtime.Spec {
	return &runtime.Spec{
		Env: []string{"SEQUENCE=true"},
	}
}

var (
	workdir = ""
)

func init() {
	var err error
	workdir, err = os.UserHomeDir()
	if err != nil {
		workdir, _ = os.Getwd()
	}

	workdir = filepath.Join(workdir, ".sqnc")
}

func PlanStep(ctx context.Context, s *sequence.Step, opts ...PlanOpt) (*runtime.Spec, error) {
	popts := &planOpts{
		path:     ".",
		workflow: &sequence.Workflow{},
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
		ghenv, err := actions.NewEnvFromPath(popts.path)
		if err != nil {
			return nil, err
		}

		action, err := actions.ParseReference(s.Uses)
		if err != nil {
			return nil, err
		}

		var (
			// generate a unique, reproducible, directory-name-compliant ID from the current context
			// so that steps that are a part of the same job share the same mounts
			id = base64.URLEncoding.EncodeToString(
				sha1.New().Sum(
					[]byte(popts.path + popts.workflow.Name + popts.jobName),
				),
			)
			gitHubEnv  = filepath.Join(workdir, id, "github", "env")
			gitHubPath = filepath.Join(workdir, id, "github", "path")
		)
		spec.Mounts = append(spec.Mounts, []specs.Mount{
			{
				Source:      filepath.Join(workdir, id, "workspace"),
				Destination: ghenv.Workspace,
				Type:        runtime.MountTypeBind,
			},
			{
				// actions can be global since every step that uses actions/checkout@v2
				// expects to function the same
				Source:      filepath.Join(workdir, "actions", action.Owner(), action.Repository(), action.Path(), action.Version()),
				Destination: ghenv.ActionPath,
				Type:        runtime.MountTypeBind,
			},
			{
				Source:      filepath.Join(workdir, id, "runner", "temp"),
				Destination: ghenv.RunnerTemp,
				Type:        runtime.MountTypeBind,
			},
			{
				Source:      filepath.Join(workdir, id, "runner", "toolcache"),
				Destination: ghenv.RunnerToolCache,
				Type:        runtime.MountTypeBind,
			},
			// these are _files_, NOT directories, that are used like
			// $ echo "MY_VAR=myval" >> $GITHUB_ENV
			// $ echo "/.mybin" >> $GITHUB_PATH
			// respectively. TODO source the contents of these files into spec.Env
			{
				Source:      gitHubEnv,
				Destination: ghenv.Env,
				Type:        runtime.MountTypeBind,
			},
			{
				Source:      gitHubPath,
				Destination: ghenv.Path,
				Type:        runtime.MountTypeBind,
			},
		}...)
		e := ghenv.Arr()
		e = append(e, env.MapToArr(s.Env)...)
		spec.Env = append(spec.Env, e...)
		spec.Cwd = ghenv.Workspace

		// s.Run doesn't necessarily need this image the way s.Uses does, but we may as well use it
		// since we own it and users will likely already have it stored locally
		spec.Image = meta.Image()

		if s.Uses != "" {
			spec.Cmd = []string{"plugin", "uses", s.Uses, ghenv.ActionPath}
			inputs := []string{}
			for input, val := range s.With {
				inputs = append(inputs, strings.ToUpper(strings.ReplaceAll(input, "-", "_")), val)
			}
			spec.Env = append(spec.Env, env.ToArr(inputs...)...)
		} else if s.Run != "" {
			switch s.Shell {
			case "/bin/bash", "bash":
				spec.Entrypoint = []string{"/bin/bash", "-c"}
			case "/bin/sh", "sh", "":
				spec.Entrypoint = []string{"/bin/sh", "-c"}
			default:
				return nil, fmt.Errorf("unsupported shell %s", s.Shell)
			}
			spec.Cmd = []string{s.Run}
		}
	} else {
		spec.Image = s.Image
		spec.Entrypoint = s.Entrypoint
		spec.Cmd = s.Cmd
	}

	return spec, nil
}