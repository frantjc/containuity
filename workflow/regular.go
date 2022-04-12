package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/env"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/net/context"
)

type regularStep struct {
	ID    string
	Name  string
	Env   map[string]string
	Shell string
	Run   string
	If    interface{}
	With  map[string]string

	Image      string
	Entrypoint []string
	Cmd        []string
	Privileged bool

	stateKey string
	specOpts []runtime.SpecOpt
}

var _ executable = &regularStep{}

func (e *regularStep) id() string {
	if e.ID != "" {
		return e.ID
	}

	return e.Name
}

func (e *regularStep) execute(ctx context.Context, ex *jobExecutor) error {
	var (
		logout   = log.New(ex.stdout).SetVerbose(ex.verbose)
		expanded = &regularStep{
			ID:         ex.expandString(e.ID),
			Name:       ex.expandString(e.Name),
			Shell:      ex.expandString(e.Shell),
			Run:        ex.expandString(e.Run),
			If:         ex.expandString(fmt.Sprint(e.If)),
			With:       ex.expandStringMap(e.With),
			Image:      ex.expandString(e.Image),
			Privileged: e.Privileged,
			Env:        ex.expandStringMap(e.Env),
			Entrypoint: ex.expandStringArr(e.Entrypoint),
			Cmd:        ex.expandStringArr(e.Cmd),
		}
		githubEnv  = ex.githubEnvFilepath()
		githubPath = ex.githubPathFilepath()
		id         = expanded.Name
		spec       = &runtime.Spec{
			Image: ex.runnerImage,
			Cwd:   ex.globalContext.GitHubContext.Workspace,
			Env: append(
				ex.globalContext.EnvArr(),
				"SQNC=true",
				"SEQUENCE=true",
				"DEBIAN_FRONTEND=noninteractive",
			),
			Mounts: ex.dirMounts(),
		}
	)

	if expanded.ID != "" {
		id = expanded.ID
	}
	ex.globalContext.InputsContext = e.With
	for k, v := range e.Env {
		ex.globalContext.EnvContext[k] = v
	}
	ex.globalContext.StepsContext[id] = &actions.StepsContext{
		Outputs: map[string]string{},
	}

	logout.Infof("[%sSQNC%s] running step '%s'", log.ColorInfo, log.ColorNone, id)

	var (
		echo                  = false
		stopCommandsTokens    = map[string]bool{}
		commandWriterCallback = func(c *actions.Command) []byte {
			if _, ok := stopCommandsTokens[c.Command]; ok {
				stopCommandsTokens[c.Command] = false
				if ex.verbose {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s end token '%s'", log.ColorDebug, log.ColorNone, actions.CommandStopCommands, c.Command))
				}
				return make([]byte, 0)
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
				if ex.verbose || echo {
					return []byte(fmt.Sprintf("[%sACTN:DBG%s] %s", log.ColorDebug, log.ColorNone, c.Value))
				}
			case actions.CommandSetOutput:
				ex.globalContext.StepsContext[id].Outputs[c.Parameters["name"]] = c.Value
				if ex.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s for '%s'", log.ColorDebug, log.ColorNone, c.Command, c.Parameters["name"], c.Value, id))
				}
			case actions.CommandStopCommands:
				stopCommandsTokens[c.Value] = true
				if ex.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s until '%s'", log.ColorDebug, log.ColorNone, c.Command, c.Value))
				}
			case actions.CommandEcho:
				if c.Value == "on" {
					echo = true
				} else if c.Value == "off" {
					echo = false
				}
			case actions.CommandSaveState:
				if e.stateKey != "" {
					ex.states[e.stateKey][c.Parameters["name"]] = c.Value
					if ex.verbose || echo {
						return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s", log.ColorDebug, log.ColorNone, c.Command, c.Parameters["name"], c.Value))
					}
				} else if ex.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] swallowing unrecognized workflow command '%s'", log.ColorDebug, log.ColorNone, c.Command))
				}
			default:
				if ex.verbose || echo {
					return []byte(fmt.Sprintf("[%sSQNC:DBG%s] swallowing unrecognized workflow command '%s'", log.ColorDebug, log.ColorNone, c.Command))
				}
			}
			return make([]byte, 0)
		}
	)

	if e.stateKey != "" {
		for k, v := range ex.states[e.stateKey] {
			spec.Env = append(spec.Env, fmt.Sprintf("STATE_%s=%s", k, v))
		}
	}

	if expanded.With != nil {
		for k, v := range expanded.With {
			spec.Env = append(spec.Env, fmt.Sprintf("INPUT_%s=%s", strings.ReplaceAll(strings.ToUpper(k), " ", "_"), v))
		}
	}

	githubEnvFile, err := createOrOpen(githubEnv)
	if err != nil {
		return err
	}

	if githubEnvArr, err := env.ArrFromReader(githubEnvFile); err == nil {
		spec.Env = append(spec.Env, githubEnvArr...)
	}

	githubPathFile, err := createOrOpen(githubPath)
	if err != nil {
		return err
	}

	if githubPathStr, err := env.PathFromReader(githubPathFile); err == nil && githubPathStr != "" {
		// TODO this overrides the default PATH instead of adding to it
		// spec.Env = append(spec.Env, env.ToArr("PATH", githubPathStr)...)
	}

	if expanded.Image != "" {
		spec.Image = expanded.Image
	}

	if expanded.Run != "" {
		switch expanded.Shell {
		case "/bin/bash", "bash":
			spec.Entrypoint = []string{"/bin/bash", "-c", expanded.Run}
		case "/bin/sh", "sh", "":
			spec.Entrypoint = []string{"/bin/sh", "-c", expanded.Run}
		default:
			return fmt.Errorf("unsupported shell '%s'", expanded.Shell)
		}
	} else {
		spec.Entrypoint = expanded.Entrypoint
		spec.Cmd = expanded.Cmd
	}

	// make sure all of the host directories that we intend to bind exist
	// note that at this point all bind mounts are directories
	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			if err = os.MkdirAll(mount.Source, 0777); err != nil {
				return err
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
			Source:      hostsFile,
			Destination: hostsFile,
			Type:        runtime.MountTypeBind,
			Options:     readOnly,
		},
		{
			Source:      resolveConf,
			Destination: resolveConf,
			Type:        runtime.MountTypeBind,
			Options:     readOnly,
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

	logout.Infof("[%sSQNC%s] pulling image '%s'", log.ColorInfo, log.ColorNone, spec.Image)
	image, err := ex.runtime.PullImage(ctx, spec.Image)
	if err != nil {
		return err
	}
	logout.Debugf("[%sSQNC:DBG%s] finished pulling image '%s'", log.ColorDebug, log.ColorNone, image.Ref())

	if e.specOpts != nil {
		for _, opt := range e.specOpts {
			if err = opt(spec); err != nil {
				return err
			}
		}
	}

	container, err := ex.runtime.CreateContainer(ctx, spec)
	if err != nil {
		return err
	}

	return container.Exec(
		ctx,
		runtime.ExecStreams(
			os.Stdin,
			actions.NewCommandWriter(commandWriterCallback, ex.stdout),
			actions.NewCommandWriter(commandWriterCallback, ex.stderr),
		),
	)
}
