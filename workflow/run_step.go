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
		ghctx            *actions.GlobalContext
		err              error
		logout           = log.New(ro.stdout)
		logerr           = log.New(ro.stderr)
	)
	logout.SetVerbose(ro.verbose)
	logerr.SetVerbose(ro.verbose)
	if ghctx, err = actions.Context(ctx); err != nil {
		copts := []actions.CtxOpt{
			actions.WithToken(ro.githubToken),
			actions.WithSecrets(ro.secrets),
			actions.WithWorkdir(containerWorkdir),
			actions.WithJobName(ro.jobName),
		}
		if ghctx, err = actions.NewContextFromPath(ctx, ro.repository, copts...); err != nil {
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
		echo                  = false
		stopCommandsTokens    = map[string]bool{}
		commandWriterCallback = func(c *actions.Command) []byte {
			if _, ok := stopCommandsTokens[c.Command]; ok {
				stopCommandsTokens[c.Command] = false
				if ro.verbose {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s end token '%s'", log.ColorDebug, log.ColorNone, actions.CommandStopCommands, c.Command))
				} else {
					return make([]byte, 0)
				}
			}

			for _, stop := range stopCommandsTokens {
				if stop {
					return []byte(c.String())
				}
			}

			switch c.Command {
			case actions.CommandError:
				return []byte(fmt.Sprintf("[%sACTN:ERR%s] %s", log.ColorError, log.ColorNone, c.Value))
			case actions.CommandWarning:
				return []byte(fmt.Sprintf("[%sACTN:WRN%s] %s", log.ColorWarn, log.ColorNone, c.Value))
			case actions.CommandNotice:
				return []byte(fmt.Sprintf("[%sACTN:NTC%s] %s", log.ColorNotice, log.ColorNone, c.Value))
			case actions.CommandDebug:
				if ro.verbose || echo {
					return []byte(fmt.Sprintf("[%sACTN:DBG%s] %s", log.ColorDebug, log.ColorNone, c.Value))
				}
			case actions.CommandSetOutput:
				if stepContext, ok := ghctx.StepsContext[es.GetID()]; ok {
					if stepContext.Outputs == nil {
						stepContext.Outputs = map[string]string{}
					}
					stepContext.Outputs[c.Parameters["name"]] = c.Value
				} else {
					ghctx.StepsContext[es.GetID()] = &actions.StepsContext{
						Outputs: map[string]string{
							c.Parameters["name"]: c.Value,
						},
					}
				}
				if ro.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s for '%s'", log.ColorDebug, log.ColorNone, c.Command, c.Parameters["name"], c.Value, es.GetID()))
				}
			case actions.CommandSaveState:
				spec.Env = append(spec.Env, fmt.Sprintf("STATE_%s=%s", c.Parameters["name"], c.Value))
				if ro.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s", log.ColorDebug, log.ColorNone, c.Command, c.Parameters["name"], c.Value))
				}
			case actions.CommandStopCommands:
				stopCommandsTokens[c.Value] = true
				if ro.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s until '%s'", log.ColorDebug, log.ColorNone, c.Command, c.Value))
				}
			case actions.CommandEcho:
				if c.Value == "on" {
					echo = true
				} else if c.Value == "off" {
					echo = false
				}
			default:
				if ro.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] swallowing unrecognized workflow command '%s'", log.ColorDebug, log.ColorNone, c.Command))
				}
			}
			return make([]byte, 0)
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

	if strings.HasPrefix(es.Uses, imagePrefix) {
		// handle uses: docker://some/action:latest
		// https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#example-using-a-docker-hub-action
		spec.Image = strings.TrimPrefix(es.Uses, imagePrefix)
	} else if es.Uses != "" {
		// handle uses: actions/checkout@v2
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
		spec.Env = append(
			spec.Env,
			fmt.Sprintf("%s=%s/%s", actions.EnvVarActionRepository, action.Owner(), action.Repository()),
		)
	} else {
		if es.Image != "" {
			spec.Image = es.Image
		}

		if es.Run != "" {
			switch es.Shell {
			case "/bin/bash", "bash":
				spec.Entrypoint = []string{"/bin/bash", "-c", es.Run}
			case "/bin/sh", "sh", "":
				spec.Entrypoint = []string{"/bin/sh", "-c", es.Run}
			default:
				return ctx, nil, fmt.Errorf("unsupported shell '%s'", es.Shell)
			}
		} else {
			spec.Entrypoint = es.Entrypoint
			spec.Cmd = es.Cmd
		}
	}

	// make sure all of the host directories that we intend to bind exist
	// note that at this point all bind mounts are directories
	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			if err = os.MkdirAll(mount.Source, 0777); err != nil {
				return ctx, nil, err
			}
		}
	}

	var (
		ghEnv  = filepath.Join(containerWorkdir, "github", "env")
		ghPath = filepath.Join(containerWorkdir, "github", "path")
	)
	spec.Env = append(
		spec.Env,
		fmt.Sprintf("%s=%s", actions.EnvVarEnv, ghEnv),
		fmt.Sprintf("%s=%s", actions.EnvVarPath, ghPath),
	)
	// these are _files_, NOT directories
	// now that we have done all of the set up for the directories we
	// intend to bind, we can add the files we intend to bind
	spec.Mounts = append(spec.Mounts, []specs.Mount{
		// make networking stuff act more predictably for users
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
		errbuf = actions.NewCommandWriter(commandWriterCallback, ro.stderr)
		eopts  = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, outbuf, errbuf)}
	)
	if !s.IsStdoutResponse() {
		eopts[0] = runtime.WithStreams(os.Stdin, actions.NewCommandWriter(commandWriterCallback, ro.stdout), errbuf)
	}
	logout.Infof("[%sSQNC%s] running step '%s'", log.ColorInfo, log.ColorNone, s.GetID())
	if err = runSpec(ctx, r, spec, ro, logout, logerr, eopts); err != nil {
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

			steps, err := NewStepsFromMetadata(action, ghctx.GitHubContext.ActionPath)
			if err != nil {
				return ctx, nil, err
			}

			var (
				outbuf = actions.NewCommandWriter(commandWriterCallback, logout)
				eopts  = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, outbuf, errbuf)}
			)
			for _, step := range steps {
				es, err = expandStep(es.Merge(&step).Canonical(), ghctx)
				if err != nil {
					return ctx, nil, err
				}

				// TODO composite actions can contain other actions,
				//      so should we recurse for composite actions?
				// if es.Uses != "" {
				// 	if ctx, _, err := runStep(ctx, r, es, ro); err != nil {
				// 		return ctx, nil, err
				// 	}
				// }

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

				logout.Infof("[%sSQNC%s] running action '%s'", log.ColorInfo, log.ColorNone, s.Uses)
				err = runSpec(ctx, r, spec, ro, logout, logerr, eopts)
				if err != nil {
					return ctx, nil, err
				}
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

func runSpec(ctx context.Context, r runtime.Runtime, s *runtime.Spec, ro *runOpts, logout, logerr log.Logger, opts []runtime.ExecOpt) error {
	logout.Infof("[%sSQNC%s] pulling image '%s'", log.ColorInfo, log.ColorNone, s.Image)
	image, err := r.PullImage(ctx, s.Image)
	if err != nil {
		return err
	}
	logout.Debugf("[%sSQNC:DBG%s] finished pulling image '%s'", log.ColorDebug, log.ColorNone, image.Ref())

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
	if err := os.MkdirAll(filepath.Dir(name), 0777); err != nil {
		return nil, err
	}

	if fs, err := os.Stat(name); err == nil && !fs.IsDir() {
		return os.Open(name)
	}

	return os.Create(name)
}
