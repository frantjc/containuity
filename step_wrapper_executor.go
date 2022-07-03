package sequence

import (
	"context"
	"fmt"
	"path"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/github/actions"

	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/internal/shim"
	"github.com/frantjc/sequence/runtime"
)

type stepWrapperExecutor struct {
	*executor
	stepWrapper        *stepWrapper
	echo               bool
	stopCommandsTokens map[string]bool
}

func (e *stepWrapperExecutor) WorkflowCommandWriterCallback(wc *actions.WorkflowCommand) []byte {
	if _, ok := e.stopCommandsTokens[wc.Command]; ok {
		e.OnWorkflowCommand.Hook(wc)
		e.stopCommandsTokens[wc.Command] = false
		if e.Verbose {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s end token '%s'", log.ColorDebug, log.ColorNone, actions.CommandStopCommands, wc.Command))
		}
		return make([]byte, 0)
	}

	for _, stop := range e.stopCommandsTokens {
		if stop {
			return []byte(wc.String())
		}
	}

	e.OnWorkflowCommand.Hook(wc)

	switch wc.Command {
	case actions.CommandError:
		return []byte(fmt.Sprintf("[%sACTN:ERR%s] %s", log.ColorError, log.ColorNone, wc.Value))
	case actions.CommandWarning:
		return []byte(fmt.Sprintf("[%sACTN:WRN%s] %s", log.ColorWarn, log.ColorNone, wc.Value))
	case actions.CommandNotice:
		return []byte(fmt.Sprintf("[%sACTN:NTC%s] %s", log.ColorNotice, log.ColorNone, wc.Value))
	case actions.CommandDebug:
		if e.Verbose || e.echo {
			return []byte(fmt.Sprintf("[%sACTN:DBG%s] %s", log.ColorDebug, log.ColorNone, wc.Value))
		}
	case actions.CommandSetOutput:
		e.GlobalContext.StepsContext["TODO"].Outputs[wc.Parameters["name"]] = wc.Value
		if e.Verbose || e.echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s for", log.ColorDebug, log.ColorNone, wc.Command, wc.Parameters["name"], wc.Value))
		}
	case actions.CommandStopCommands:
		e.stopCommandsTokens[wc.Value] = true
		if e.Verbose || e.echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s until '%s'", log.ColorDebug, log.ColorNone, wc.Command, wc.Value))
		}
	case actions.CommandEcho:
		if wc.Value == "on" {
			e.echo = true
		} else if wc.Value == "off" {
			e.echo = false
		}
	case actions.CommandSaveState:
		e.stepWrapper.state[wc.Parameters["name"]] = wc.Value
		if e.Verbose || e.echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s", log.ColorDebug, log.ColorNone, wc.Command, wc.Parameters["name"], wc.Value))
		}
	default:
		if e.Verbose || e.echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] swallowing unrecognized workflow command '%s'", log.ColorDebug, log.ColorNone, wc.Command))
		}
	}
	return make([]byte, 0)
}

func (e *stepWrapperExecutor) ExecuteStep(ctx context.Context) error {
	var (
		// logStdout    = log.New(e.stdout).SetVerbose(e.verbose)
		// logStderr    = log.New(e.stderr).SetVerbose(e.verbose)
		expander     = actions.NewExpander(e.GlobalContext.Get)
		expandedStep = &Step{
			Id:         expander.Expand(e.stepWrapper.step.Id),
			Name:       expander.Expand(e.stepWrapper.step.Name),
			Shell:      expander.Expand(e.stepWrapper.step.Shell),
			Run:        expander.Expand(e.stepWrapper.step.Run),
			If:         expander.Expand(e.stepWrapper.step.If),
			Image:      expander.Expand(e.stepWrapper.step.Image),
			Privileged: e.stepWrapper.step.Privileged,
			Entrypoint: js.Map(e.stepWrapper.step.Entrypoint, func(arg string, _ int, _ []string) string {
				return expander.Expand(arg)
			}),
			Cmd: js.Map(e.stepWrapper.step.Cmd, func(arg string, _ int, _ []string) string {
				return expander.Expand(arg)
			}),
			Env:  map[string]string{},
			With: map[string]string{},
		}
		id = js.Coalesce(expandedStep.Id, expandedStep.Name)
	)

	for k, v := range e.stepWrapper.step.Env {
		expandedStep.Env[k] = expander.Expand(v)
	}

	for k, v := range e.stepWrapper.step.With {
		expandedStep.With[k] = expander.Expand(v)
	}

	e.GlobalContext.InputsContext = expandedStep.With
	e.GlobalContext.AddEnv(expandedStep.Env)
	e.GlobalContext.AddEnv(e.stepWrapper.step.Env)
	e.GlobalContext.AddEnv(e.stepWrapper.extraEnv)
	e.GlobalContext.StepsContext[id] = &actions.StepsContext{
		Outputs: map[string]string{},
	}

	spec := &runtime.Spec{
		Image:      js.Coalesce(expandedStep.Image, e.RunnerImage.GetRef()),
		Entrypoint: []string{paths.Shim},
		Cwd:        e.GlobalContext.GitHubContext.Workspace,
		Env: append(
			[]string{
				"SQNC=true",
				"SEQUENCE=true",
				fmt.Sprintf("%s=", shim.EnvVarShimSwitch),
				fmt.Sprintf("%s=%s", actions.EnvVarEnv, paths.GitHubEnv),
				fmt.Sprintf("%s=%s", actions.EnvVarPath, paths.GitHubPath),
			},
			e.GlobalContext.EnvArr()...,
		),
		Mounts: append(
			[]*runtime.Mount{
				{
					Source:      volumes.GetWorkspace(e.stepWrapper.id),
					Destination: e.GlobalContext.GitHubContext.Workspace,
					Type:        runtime.MountTypeVolume,
				},
				{
					Source:      volumes.GetRunnerTmp(e.stepWrapper.id),
					Destination: e.GlobalContext.RunnerContext.Temp,
					Type:        runtime.MountTypeVolume,
				},
				{
					Source:      volumes.GetRunnerToolCache(e.stepWrapper.id),
					Destination: e.GlobalContext.RunnerContext.ToolCache,
					Type:        runtime.MountTypeVolume,
				},
				{
					// paths.GitHubEnv and paths.GitHubPath are files in the same
					// directory, so we only need one mount for both to share
					Source:      volumes.GetGitHub(e.stepWrapper.id),
					Destination: path.Dir(paths.GitHubPath),
					Type:        runtime.MountTypeVolume,
				},
			},
			e.stepWrapper.extraMounts...,
		),
	}

	for k, v := range e.stepWrapper.state {
		spec.Env = append(spec.Env, fmt.Sprintf("STATE_%s=%s", k, v))
	}

	if expandedStep.Run != "" {
		switch expandedStep.Shell {
		case "/bin/bash", "bash":
			spec.Cmd = []string{"/bin/bash", "-c", expandedStep.Run}
		case "/bin/sh", "sh", "":
			spec.Cmd = []string{"/bin/sh", "-c", expandedStep.Run}
		default:
			return fmt.Errorf("unsupported shell '%s'", expandedStep.Shell)
		}
	} else {
		spec.Cmd = append(expandedStep.Entrypoint, expandedStep.Cmd...) //nolint:gocritic
	}

	return e.RunContainer(
		ctx,
		spec,
		runtime.NewStreams(
			e.Stdin,
			actions.NewWorkflowCommandWriter(e.WorkflowCommandWriterCallback, e.Stdout),
			actions.NewWorkflowCommandWriter(e.WorkflowCommandWriterCallback, e.Stderr),
		),
	)
}
