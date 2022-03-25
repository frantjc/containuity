package workflow

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/env"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var readonly = []string{runtime.MountOptReadOnly}

func RunStep(ctx context.Context, r runtime.Runtime, s *Step, opts ...RunOpt) (context.Context, *StepOut, error) {
	ro, err := newRunOpts(opts...)
	if err != nil {
		return ctx, nil, err
	}

	return runStep(ctx, r, s, ro)
}

func runStep(ctx context.Context, r runtime.Runtime, s *Step, ro *runOpts) (context.Context, *StepOut, error) {
	var (
		containerWorkdir = "/sqnc"
		copts            = []actions.CtxOpt{
			actions.WithToken(ro.githubToken),
			actions.WithSecrets(ro.secrets),
			actions.WithWorkdir(containerWorkdir),
			actions.WithJobName(ro.jobName),
		}
		ghctx *actions.GlobalContext
	)
	ghctx, err := actions.Context(ctx)
	if err != nil {
		ghctx, err = actions.NewContextFromPath(ctx, ro.repository, copts...)
		if err != nil {
			return ctx, nil, err
		}
	}

	es, err := expandStep(s.Canonical(), ghctx)
	if err != nil {
		return ctx, nil, err
	}

	if err = actions.WithInputs(es.With)(ghctx); err != nil {
		return ctx, nil, err
	}

	var (
		// generate a unique, reproducible, directory-name-compliant ID from the current context
		// so that steps that are a part of the same job share the same mounts
		id = base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprint(ro.repository, ro.workflow.Name, ro.jobName)),
		)
		workdirid  = filepath.Join(ro.workdir, id)
		githubEnv  = filepath.Join(workdirid, "github", "env")
		githubPath = filepath.Join(workdirid, "github", "path")
		spec       = &runtime.Spec{
			Image:      ro.runnerImage,
			Cwd:        ghctx.GitHubContext.Workspace,
			Privileged: es.Privileged,
			Env: append(
				ghctx.Arr(),
				"SEQUENCE=true",
				"SQNC=true",
				"DEBIAN_FRONTEND=noninteractive",
				"ACCEPT_EULA=Y",
			),
			Mounts: []specs.Mount{
				{
					Source:      "/etc/ssl",
					Destination: "/etc/ssl",
					Type:        runtime.MountTypeBind,
					Options:     readonly,
				},
				{
					Source:      filepath.Join(workdirid, "workspace"),
					Destination: ghctx.GitHubContext.Workspace,
					Type:        runtime.MountTypeBind,
				},
				{
					Source:      filepath.Join(workdirid, "runner", "temp"),
					Destination: ghctx.RunnerContext.Temp,
					Type:        runtime.MountTypeBind,
				},
				{
					Source:      filepath.Join(workdirid, "runner", "toolcache"),
					Destination: ghctx.RunnerContext.ToolCache,
					Type:        runtime.MountTypeBind,
				},
			},
		}
	)

	githubEnvFile, err := createFile(githubEnv)
	if err != nil {
		return ctx, nil, err
	}

	if githubEnvArr, err := env.ArrFromReader(githubEnvFile); err == nil {
		spec.Env = append(spec.Env, githubEnvArr...)
	}

	githubPathFile, err := createFile(githubPath)
	if err != nil {
		return ctx, nil, err
	}

	if githubPathStr, err := env.PathFromReader(githubPathFile); err == nil && githubPathStr != "" {
		spec.Env = append(spec.Env, env.ToArr("PATH", githubPathStr)...)
	}

	// TODO handle `uses: docker://some/action:latest`
	// https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#example-using-a-docker-hub-action
	if es.Uses != "" {
		action, err := actions.ParseReference(es.Uses)
		if err != nil {
			return ctx, nil, err
		}

		spec.Image = meta.Image()
		spec.Entrypoint = []string{"sqncshim"}
		spec.Cmd = []string{"plugin", "uses", action.String(), ghctx.GitHubContext.ActionPath}
		spec.Mounts = append(spec.Mounts, specs.Mount{
			// actions can be global since every step that uses actions/checkout@v2
			// expects it to function the same
			Source:      filepath.Join(ro.workdir, "actions", action.Owner(), action.Repository(), action.Path(), action.Version()),
			Destination: ghctx.GitHubContext.ActionPath,
			Type:        runtime.MountTypeBind,
		})
		spec.Env = append(spec.Env, fmt.Sprintf("%s=%s/%s", actions.EnvVarActionRepository, action.Owner(), action.Repository()))
	} else {
		if es.Image != "" {
			spec.Image = es.Image
		}

		if es.Run != "" {
			switch es.Shell {
			case "/bin/bash", "bash":
				spec.Entrypoint = []string{"/bin/bash", "-c"}
			case "/bin/sh", "sh", "":
				spec.Entrypoint = []string{"/bin/sh", "-c"}
			default:
				return ctx, nil, fmt.Errorf("unsupported shell '%s'", es.Shell)
			}
			spec.Cmd = []string{es.Run}
		} else {
			spec.Entrypoint = es.Entrypoint
			spec.Cmd = es.Cmd
		}
	}

	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			err = os.MkdirAll(mount.Source, 0777)
			if err != nil {
				return ctx, nil, err
			}
		}
	}

	var (
		ghEnv  = filepath.Join(containerWorkdir, "github", "env")
		ghPath = filepath.Join(containerWorkdir, "github", "path")
	)
	spec.Env = append(spec.Env, fmt.Sprintf("%s=%s", actions.EnvVarEnv, ghEnv), fmt.Sprintf("%s=%s", actions.EnvVarPath, ghPath))
	// these are _files_, NOT directories
	spec.Mounts = append(spec.Mounts, []specs.Mount{
		{
			Source:      "/etc/hosts",
			Destination: "/etc/hosts",
			Type:        runtime.MountTypeBind,
			Options:     readonly,
		},
		{
			Source:      "/etc/resolv.conf",
			Destination: "/etc/resolv.conf",
			Type:        runtime.MountTypeBind,
			Options:     readonly,
		},
		// these are used like
		// $ echo "MY_VAR=myval" >> $GITHUB_ENV
		// $ echo "/.mybin" >> $GITHUB_PATH
		// respectively
		{
			Source:      githubEnv,
			Destination: ghEnv,
			Type:        runtime.MountTypeBind,
		},
		{
			Source:      githubPath,
			Destination: ghPath,
			Type:        runtime.MountTypeBind,
		},
	}...)

	var (
		outbuf = new(bytes.Buffer)
		errbuf = ro.stderr
		eopts  = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, outbuf, errbuf)}
	)
	if !s.IsStdoutResponse() {
		eopts[0] = runtime.WithStreams(os.Stdin, ro.stdout, errbuf)
	}
	ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC%s] running step '%s'\n", log.ColorInfo, log.ColorNone, s.GetID())))
	if err = runSpec(ctx, r, spec, ro, eopts); err != nil {
		return ctx, nil, err
	}

	resp := &StepOut{}
	if s.IsStdoutResponse() {
		if err = json.NewDecoder(outbuf).Decode(resp); err != nil {
			return ctx, nil, err
		}

		if actionj := resp.Metadata[ActionMetadataKey]; actionj != "" {
			action := &actions.Metadata{}
			err := json.Unmarshal([]byte(actionj), action)
			if err != nil {
				return ctx, nil, err
			}

			rs, err := NewStepFromMetadata(action, ghctx.GitHubContext.ActionPath)
			if err != nil {
				return ctx, nil, err
			}

			es, err = expandStep(es.Merge(rs).Canonical(), ghctx)
			if err != nil {
				return ctx, nil, err
			}

			spec.Image = es.Image
			if ro.actionImage != "" && (action.Runs.Using == "node12" || action.Runs.Using == "node16") {
				spec.Image = ro.actionImage
			}

			spec.Entrypoint = es.Entrypoint
			spec.Cmd = es.Cmd

			for envVar, val := range es.With {
				spec.Env = append(
					spec.Env,
					fmt.Sprintf(
						"INPUT_%s=%s",
						strings.ToUpper(strings.ReplaceAll(envVar, "-", "_")),
						val,
					),
				)
			}

			var (
				outbuf = actions.NewCommandWriter(func(c *actions.Command) []byte {
					switch c.Command {
					case actions.CommandError:
						return []byte(fmt.Sprintf("[%sACTN:ERR%s] %s", log.ColorError, log.ColorNone, c.Value))
					case actions.CommandWarning:
						return []byte(fmt.Sprintf("[%sACTN:WRN%s] %s", log.ColorWarn, log.ColorNone, c.Value))
					case actions.CommandNotice:
						return []byte(fmt.Sprintf("[%sACTN:NTC%s] %s", log.ColorNotice, log.ColorNone, c.Value))
					case actions.CommandDebug:
						if ro.verbose {
							return []byte(fmt.Sprintf("[%sACTN:DBG%s] %s", log.ColorDebug, log.ColorNone, c.Value))
						}
					}
					return make([]byte, 0)
				}, ro.stdout)
				eopts = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, outbuf, errbuf)}
			)
			ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC%s] running action '%s'\n", log.ColorInfo, log.ColorNone, s.Uses)))
			err = runSpec(ctx, r, spec, ro, eopts)
			if err != nil {
				return ctx, nil, err
			}
		}
	}

	return actions.WithContext(ctx, ghctx), resp, nil
}

func expandStep(s *Step, ctx *actions.GlobalContext) (*Step, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	es := &Step{}
	err = json.Unmarshal(
		actions.ExpandBytes(b, func(s string) string {
			return ctx.Get(s)
		}),
		es,
	)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func runSpec(ctx context.Context, r runtime.Runtime, s *runtime.Spec, ro *runOpts, opts []runtime.ExecOpt) error {
	ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC%s] pulling image '%s'\n", log.ColorInfo, log.ColorNone, s.Image)))
	image, err := r.PullImage(ctx, s.Image)
	if err != nil {
		return err
	}
	if ro.verbose {
		ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC:DBG%s] finished pulling image '%s'\n", log.ColorDebug, log.ColorNone, image.Ref())))
	}

	container, err := r.CreateContainer(ctx, s)
	if err != nil {
		return err
	}

	err = container.Exec(ctx, opts...)
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

	if fs, err := os.Stat(name); err == nil && !fs.IsDir() {
		return os.Open(name)
	}

	return os.Create(name)
}
