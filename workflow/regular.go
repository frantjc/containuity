package workflow

import (
	"fmt"
	"os"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
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
	mounts   []*runtimev1.Mount
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
			ID:   ex.expandString(e.ID),
			Name: ex.expandString(e.Name),
			With: ex.expandStringMap(e.With),
		}
		spec = &runtimev1.Spec{
			Image:      ex.runnerImage,
			Entrypoint: []string{containerShim},
			Cwd:        ex.globalContext.GitHubContext.Workspace,
			Env: append(
				ex.env(),
				"SQNC=true",
				"SEQUENCE=true",
				"DEBIAN_FRONTEND=noninteractive",
			),
			Mounts: ex.mounts(),
		}
	)

	id := expanded.Name
	if expanded.ID != "" {
		id = expanded.ID
	}
	logout.Infof("[%sSQNC%s] running step '%s'", log.ColorInfo, log.ColorNone, id)
	ex.globalContext.InputsContext = expanded.With
	expanded.Env = ex.expandStringMap(e.Env)
	if ex.globalContext.EnvContext == nil {
		ex.globalContext.EnvContext = map[string]string{}
	}
	for k, v := range expanded.Env {
		ex.globalContext.EnvContext[k] = v
	}
	expanded = &regularStep{
		ID:         expanded.ID,
		Name:       expanded.Name,
		Shell:      ex.expandString(e.Shell),
		Run:        ex.expandString(e.Run),
		If:         ex.expandString(fmt.Sprint(e.If)),
		With:       expanded.With,
		Image:      ex.expandString(e.Image),
		Privileged: e.Privileged,
		Env:        ex.expandStringMap(expanded.Env),
		Entrypoint: ex.expandStringArr(e.Entrypoint),
		Cmd:        ex.expandStringArr(e.Cmd),
	}
	ex.globalContext.StepsContext[id] = &actions.StepsContext{
		Outputs: map[string]string{},
	}
	spec.Env = append(spec.Env, ex.globalContext.EnvArr()...)

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

	if expanded.Image != "" {
		spec.Image = expanded.Image
	}

	if expanded.Run != "" {
		switch expanded.Shell {
		case "/bin/bash", "bash":
			spec.Cmd = []string{"/bin/bash", "-c", expanded.Run}
		case "/bin/sh", "sh", "":
			spec.Cmd = []string{"/bin/sh", "-c", expanded.Run}
		default:
			return fmt.Errorf("unsupported shell '%s'", expanded.Shell)
		}
	} else {
		spec.Cmd = append(expanded.Entrypoint, expanded.Cmd...) //nolint:gocritic
	}

	logout.Infof("[%sSQNC%s] pulling image '%s'", log.ColorInfo, log.ColorNone, spec.Image)
	image, err := ex.runtime.PullImage(ctx, spec.Image)
	if err != nil {
		return err
	}
	logout.Debugf("[%sSQNC:DBG%s] finished pulling image '%s'", log.ColorDebug, log.ColorNone, image.GetRef())

	spec.Mounts = append(spec.Mounts, e.mounts...)

	logout.Debugf("[%sSQNC:DBG%s] getting or creating volumes", log.ColorDebug, log.ColorNone)
	for _, mount := range spec.Mounts {
		if mount.Type == runtimev1.MountTypeVolume {
			vol, err := ex.runtime.CreateVolume(ctx, mount.Source)
			if err != nil {
				if vol, err = ex.runtime.GetVolume(ctx, mount.Source); err != nil {
					return err
				}
			}
			mount.Source = vol.GetSource()
		}
	}
	logout.Debugf("[%sSQNC:DBG%s] finished setting up volumes", log.ColorDebug, log.ColorNone)

	logout.Debugf("[%sSQNC:DBG%s] creating container", log.ColorDebug, log.ColorNone)
	container, err := ex.runtime.CreateContainer(ctx, spec)
	if err != nil {
		return err
	}

	logout.Debugf("[%sSQNC:DBG%s] copying shim to container", log.ColorDebug, log.ColorNone)
	sqncshim, err := shimSourceTarArchive()
	if err != nil {
		return err
	}

	if err = container.CopyTo(ctx, sqncshim, containerShimDir); err != nil {
		return err
	}

	return container.Exec(
		ctx,
		runtime.NewStreams(
			os.Stdin,
			actions.NewCommandWriter(commandWriterCallback, ex.stdout),
			actions.NewCommandWriter(commandWriterCallback, ex.stderr),
		),
	)
}
