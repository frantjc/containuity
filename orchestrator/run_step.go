package orchestrator

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var (
	readonly = []string{runtime.MountOptReadOnly}
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

func RunStep(ctx context.Context, r runtime.Runtime, s *sequence.Step, opts ...OrchOpt) error {
	var (
		oo = &orchOpts{
			workflow: &sequence.Workflow{},
			job: &sequence.Job{},
			path: ".",
		}
	)
	for _, opt := range opts {
		err := opt(oo)
		if err != nil {
			return err
		}
	}

	ghvars, err := actions.NewVarsFromPath(oo.path)
	if err != nil {
		return err
	}

	var (
		ghenv = ghvars.Env
		ghctx = ghvars.ActionsContext
	)
	es, err := expandStep(s, ghctx)
	if err != nil {
		return err
	}

	var (
		// generate a unique, reproducible, directory-name-compliant ID from the current context
		// so that steps that are a part of the same job share the same mounts
		id = base64.URLEncoding.EncodeToString(
			sha1.New().Sum(
				[]byte(oo.path + oo.workflow.Name + oo.jobName),
			),
		)
		gitHubEnv  = filepath.Join(workdir, id, "github", "env")
		gitHubPath = filepath.Join(workdir, id, "github", "path")
		spec       = &runtime.Spec{
			Image:      meta.Image(),
			Cwd:        ghenv.Workspace,
			Privileged: es.Privileged,
			Env:        append(ghenv.Arr(), "SEQUENCE=true"),
			Mounts: []specs.Mount{
				{
					Source: "/etc/ssl",
					Destination: "/etc/ssl",
					Type: runtime.MountTypeBind,
					Options: readonly,
				},
				{
					Source:      filepath.Join(workdir, id, "workspace"),
					Destination: ghenv.Workspace,
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
			},
		}
	)

	_, err = createFile(gitHubEnv)
	if err != nil {
		return err
	}

	_, err = createFile(gitHubPath)
	if err != nil {
		return err
	}

	if es.Uses != "" {
		action, err := actions.ParseReference(es.Uses)
		if err != nil {
			return err
		}

		spec.Cmd = []string{"plugin", "uses", action.String(), ghenv.ActionPath}
		spec.Mounts = append(spec.Mounts, specs.Mount{
			// actions can be global since every step that uses actions/checkout@v2
			// expects to function the same
			Source:      filepath.Join(workdir, "actions", action.Owner(), action.Repository(), action.Path(), action.Version()),
			Destination: ghenv.ActionPath,
			Type:        runtime.MountTypeBind,
		})
	} else if es.Run != "" {
		switch es.Shell {
		case "/bin/bash", "bash":
			spec.Entrypoint = []string{"/bin/bash", "-c"}
		case "/bin/sh", "sh", "":
			spec.Entrypoint = []string{"/bin/sh", "-c"}
		default:
			return fmt.Errorf("unsupported shell '%s'", es.Shell)
		}
		spec.Cmd = []string{es.Run}
	} else if es.Image != "" {
		spec.Image = es.Image
		spec.Entrypoint = es.Entrypoint
		spec.Cmd = es.Cmd
	}

	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			err = os.MkdirAll(mount.Source, 0777)
			if err != nil {
				return err
			}
		}
	}

	spec.Mounts = append(spec.Mounts, []specs.Mount{
		{
			Source: "/etc/hosts",
			Destination: "/etc/hosts",
			Type: runtime.MountTypeBind,
			Options: readonly,
		},
		{
			Source: "/etc/resolv.conf",
			Destination: "/etc/resolv.conf",
			Type: runtime.MountTypeBind,
			Options: readonly,
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

	var (
		copts = append([]runtime.SpecOpt{runtime.WithSpec(spec)}, oo.sopts...)
		eopts = []runtime.ExecOpt{runtime.WithStdio}
		buf   = new(bytes.Buffer)
	)
	if s.IsStdoutResponse() {
		eopts[0] = runtime.WithStreams(os.Stdin, buf, os.Stderr)
	}
	err = runSpec(ctx, r, spec, copts, eopts)
	if err != nil {
		return err
	}

	if s.IsStdoutResponse() {
		resp := &sequence.StepResponse{}
		if err = json.NewDecoder(buf).Decode(resp); err != nil {
			return err
		}

		if resp.Step != nil {
			es, err = expandStep(es.Merge(resp.Step).Canonical(), ghctx)
			if err != nil {
				return err
			}

			spec.Image = es.Image
			spec.Entrypoint = es.Entrypoint
			spec.Cmd = es.Cmd

			for envvar, val := range es.With {
				spec.Env = append(
					spec.Env,
					fmt.Sprintf(
						"INPUT_%s=%s",
						strings.ToUpper(strings.ReplaceAll(envvar, "-", "_")),
						val,
					),
				)
			}

			var (
				copts = append([]runtime.SpecOpt{runtime.WithSpec(spec)}, oo.sopts...)
				eopts = []runtime.ExecOpt{runtime.WithStdio}
			)
			err = runSpec(ctx, r, spec, copts, eopts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func expandStep(s *sequence.Step, ctx *actions.ActionsContext) (*sequence.Step, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	es := &sequence.Step{}
	err = json.Unmarshal(
		actions.ExpandBytes(b, func(s string) string {
			return ctx.Value(s).(string)
		}),
		es,
	)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func runSpec(ctx context.Context, r runtime.Runtime, s *runtime.Spec, copts []runtime.SpecOpt, eopts []runtime.ExecOpt) error {
	container, err := r.Create(ctx, copts...)
	if err != nil {
		return err
	}

	err = container.Exec(ctx, eopts...)
	if err != nil {
		return err
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(name), 0777)
	if err != nil {
		return nil, err
	}

	return os.Create(name)
}
