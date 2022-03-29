package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/env"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var readonly = []string{runtime.MountOptReadOnly}

func RunStep(ctx context.Context, r runtime.Runtime, s *Step, opts ...RunOpt) error {
	ro, err := newRunOpts(opts...)
	if err != nil {
		return err
	}

	var (
		_       *StepOut
		ctxopts = []actions.CtxOpt{
			actions.WithToken(ro.githubToken),
			actions.WithSecrets(ro.secrets),
			actions.WithWorkdir(containerWorkdir),
			actions.WithJobName(ro.jobName),
		}
	)
	ro.gctx, err = actions.NewContextFromPath(ctx, ro.repository, ctxopts...)
	if err != nil {
		return err
	}

	if s.Uses != "" {
		expandedStep, err := expandStep(s.Canonical(), ro.gctx)
		if err != nil {
			return err
		}

		if metadata, reference, ro, err := runStepSetup(ctx, r, expandedStep, ro); err != nil {
			return err
		} else if metadata != nil && reference != nil {
			ro.specOpts = append(ro.specOpts, runtime.WithMounts(
				specs.Mount{
					Source:      getHostActionPath(reference, ro),
					Destination: ro.gctx.GitHubContext.ActionPath,
					Type:        runtime.MountTypeBind,
				},
			))
			with := withFromInputs(metadata.Inputs)
			if preStep, err := newPreStepFromMetadataWith(metadata, ro.gctx.GitHubContext.ActionPath, with); err != nil {
				return err
			} else if preStep != nil {
				if ro.actionImage != "" {
					preStep.Image = ro.actionImage
				}
				if expandedStep.ID != "" {
					preStep.ID = fmt.Sprintf("Pre %s", expandedStep.ID)
				}
				if expandedStep.Name != "" {
					preStep.Name = fmt.Sprintf("Pre %s", expandedStep.Name)
				}

				expandedPreStep, err := expandStep(preStep.Canonical(), ro.gctx)
				if err != nil {
					return err
				}

				ro.logout.Debugf("[%sSQNC:DBG%s] running action pre step '%s'", log.ColorDebug, log.ColorNone, expandedPreStep.GetID())
				if _, ro, err = runStep(ctx, r, expandedPreStep, ro); err != nil {
					return err
				}
			}

			if mainStep, err := newMainStepFromMetadataWith(metadata, ro.gctx.GitHubContext.ActionPath, with); err != nil {
				return err
			} else if mainStep != nil {
				if ro.actionImage != "" {
					mainStep.Image = ro.actionImage
				}
				mainStep.ID = expandedStep.ID
				mainStep.Name = expandedStep.Name

				expandedMainStep, err := expandStep(mainStep.Canonical(), ro.gctx)
				if err != nil {
					return err
				}

				ro.logout.Debugf("[%sSQNC:DBG%s] running action main step '%s'", log.ColorDebug, log.ColorNone, expandedMainStep.GetID())
				if _, ro, err = runStep(ctx, r, expandedMainStep, ro); err != nil {
					return err
				}
			}

			if postStep, err := newPostStepFromMetadataWith(metadata, ro.gctx.GitHubContext.ActionPath, with); err != nil {
				return err
			} else if postStep != nil {
				if ro.actionImage != "" {
					postStep.Image = ro.actionImage
				}
				if expandedStep.ID != "" {
					postStep.ID = fmt.Sprintf("Post %s", expandedStep.ID)
				}
				if expandedStep.Name != "" {
					postStep.Name = fmt.Sprintf("Post %s", expandedStep.Name)
				}

				expandedPostStep, err := expandStep(postStep.Canonical(), ro.gctx)
				if err != nil {
					return err
				}

				ro.logout.Debugf("[%sSQNC:DBG%s] running action post step '%s'", log.ColorDebug, log.ColorNone, expandedPostStep.GetID())
				if _, ro, err = runStep(ctx, r, expandedPostStep, ro); err != nil {
					return err
				}
			}
		}
	} else {
		expandedStep, err := expandStep(s.Canonical(), ro.gctx)
		if err != nil {
			return err
		}

		ro.logout.Debugf("[%sSQNC:DBG%s] running raw step '%s'", log.ColorDebug, log.ColorNone, expandedStep.GetID())
		if _, ro, err = runStep(ctx, r, expandedStep, ro); err != nil {
			return err
		}
	}

	return nil
}

func runStep(ctx context.Context, r runtime.Runtime, expandedStep *Step, ro *runOpts) (*StepOut, *runOpts, error) {
	if expandedStep.Uses != "" {
		return nil, nil, fmt.Errorf("steps with uses must be setup and shed their uses before being ran")
	}

	var (
		// generate a unique, reproducible, directory-name-compliant ID from the current context
		// so that steps that are a part of the same job share the same mounts
		stepID                = getID(ro)
		stepWorkdir           = getHostWorkdir(stepID, ro)
		githubEnv             = getHostGitHubEnvFilepath(stepWorkdir)
		githubPath            = getHostGitHubPathFilepath(stepWorkdir)
		spec                  = getDefaultSpec(ro.gctx, stepWorkdir, expandedStep.Privileged, ro)
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
				if stepContext, ok := ro.gctx.StepsContext[expandedStep.GetID()]; ok {
					if stepContext.Outputs == nil {
						stepContext.Outputs = map[string]string{}
					}
					stepContext.Outputs[c.Parameters["name"]] = c.Value
				} else {
					ro.gctx.StepsContext[expandedStep.GetID()] = &actions.StepsContext{
						Outputs: map[string]string{
							c.Parameters["name"]: c.Value,
						},
					}
				}
				if ro.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s for '%s'", log.ColorDebug, log.ColorNone, c.Command, c.Parameters["name"], c.Value, expandedStep.GetID()))
				}
			case actions.CommandSaveState:
				ro.specOpts = append(
					ro.specOpts,
					runtime.WithEnv(map[string]string{
						fmt.Sprintf("STATE_%s", c.Parameters["name"]): c.Value,
					}),
				)
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
		return nil, ro, err
	}

	if githubEnvArr, err := env.ArrFromReader(githubEnvFile); err == nil {
		spec.Env = append(spec.Env, githubEnvArr...)
	}

	githubPathFile, err := createFile(githubPath)
	if err != nil {
		return nil, ro, err
	}

	if githubPathStr, err := env.PathFromReader(githubPathFile); err == nil && githubPathStr != "" {
		// TODO this overrides the default PATH instead of adding to it
		// spec.Env = append(spec.Env, env.ToArr("PATH", githubPathStr)...)
	}

	if strings.HasPrefix(expandedStep.Uses, imagePrefix) {
		// handle `uses: docker://some/action:latest`
		// https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#example-using-a-docker-hub-action
		spec.Image = strings.TrimPrefix(expandedStep.Uses, imagePrefix)
	} else {
		if expandedStep.Image != "" {
			spec.Image = expandedStep.Image
		}

		if expandedStep.Run != "" {
			switch expandedStep.Shell {
			case "/bin/bash", "bash":
				spec.Entrypoint = []string{"/bin/bash", "-c", expandedStep.Run}
			case "/bin/sh", "sh", "":
				spec.Entrypoint = []string{"/bin/sh", "-c", expandedStep.Run}
			default:
				return nil, ro, fmt.Errorf("unsupported shell '%s'", expandedStep.Shell)
			}
		} else {
			spec.Entrypoint = expandedStep.Entrypoint
			spec.Cmd = expandedStep.Cmd
		}
	}

	// make sure all of the host directories that we intend to bind exist
	// note that at this point all bind mounts are directories
	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			if err = os.MkdirAll(mount.Source, 0777); err != nil {
				return nil, ro, err
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
	if !expandedStep.IsStdoutResponse() {
		eopts[0] = runtime.WithStreams(os.Stdin, actions.NewCommandWriter(commandWriterCallback, ro.stdout), errbuf)
	}
	if err = runSpec(ctx, r, spec, ro, eopts); err != nil {
		return nil, ro, err
	}

	resp := &StepOut{}
	if expandedStep.IsStdoutResponse() {
		if err = json.NewDecoder(outbuf).Decode(resp); err != nil {
			return nil, ro, err
		}
	}

	return resp, ro, nil
}
