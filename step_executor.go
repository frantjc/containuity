package sequence

import (
	"context"
	"fmt"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/github/actions"

	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/internal/shim"
	"github.com/frantjc/sequence/runtime"
)

type stepExecutor struct {
	executor
	StepWrapper        *StepWrapper
	Echo               bool
	StopCommandsTokens map[string]bool
}

func (e *stepExecutor) WorkflowCommandWriterCallback(wc *actions.WorkflowCommand) []byte {
	e.OnWorkflowCommand.Hook(wc)

	if _, ok := e.StopCommandsTokens[wc.Command]; ok {
		e.StopCommandsTokens[wc.Command] = false
		if e.Verbose {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s end token '%s'", log.ColorDebug, log.ColorNone, actions.CommandStopCommands, wc.Command))
		}
		return make([]byte, 0)
	}

	for _, stop := range e.StopCommandsTokens {
		if stop {
			return []byte(wc.String())
		}
	}

	switch wc.Command {
	case actions.CommandError:
		return []byte(fmt.Sprintf("[%sACTN:ERR%s] %s", log.ColorError, log.ColorNone, wc.Value))
	case actions.CommandWarning:
		return []byte(fmt.Sprintf("[%sACTN:WRN%s] %s", log.ColorWarn, log.ColorNone, wc.Value))
	case actions.CommandNotice:
		return []byte(fmt.Sprintf("[%sACTN:NTC%s] %s", log.ColorNotice, log.ColorNone, wc.Value))
	case actions.CommandDebug:
		if e.Verbose || e.Echo {
			return []byte(fmt.Sprintf("[%sACTN:DBG%s] %s", log.ColorDebug, log.ColorNone, wc.Value))
		}
	case actions.CommandSetOutput:
		e.GlobalContext.StepsContext["TODO"].Outputs[wc.Parameters["name"]] = wc.Value
		if e.Verbose || e.Echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s for", log.ColorDebug, log.ColorNone, wc.Command, wc.Parameters["name"], wc.Value))
		}
	case actions.CommandStopCommands:
		e.StopCommandsTokens[wc.Value] = true
		if e.Verbose || e.Echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s until '%s'", log.ColorDebug, log.ColorNone, wc.Command, wc.Value))
		}
	case actions.CommandEcho:
		if wc.Value == "on" {
			e.Echo = true
		} else if wc.Value == "off" {
			e.Echo = false
		}
	case actions.CommandSaveState:
		e.StepWrapper.State[wc.Parameters["name"]] = wc.Value
		if e.Verbose || e.Echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] %s %s=%s", log.ColorDebug, log.ColorNone, wc.Command, wc.Parameters["name"], wc.Value))
		}
	default:
		if e.Verbose || e.Echo {
			return []byte(fmt.Sprintf("[%sSQNC:DBG%s] swallowing unrecognized workflow command '%s'", log.ColorDebug, log.ColorNone, wc.Command))
		}
	}
	return make([]byte, 0)
}

func (e *stepExecutor) ExecuteStep(ctx context.Context) error {
	var (
		// logStdout    = log.New(e.stdout).SetVerbose(e.verbose)
		// logStderr    = log.New(e.stderr).SetVerbose(e.verbose)
		expander     = actions.NewExpander(e.GlobalContext.Get)
		expandedStep = &Step{
			Id:         expander.Expand(e.StepWrapper.Id),
			Name:       expander.Expand(e.StepWrapper.Name),
			Shell:      expander.Expand(e.StepWrapper.Shell),
			Run:        expander.Expand(e.StepWrapper.Run),
			If:         expander.Expand(e.StepWrapper.If),
			Image:      expander.Expand(e.StepWrapper.Image),
			Privileged: e.StepWrapper.Privileged,
			Entrypoint: js.Map(e.StepWrapper.Entrypoint, func(arg string, _ int, _ []string) string {
				return expander.Expand(arg)
			}),
			Cmd: js.Map(e.StepWrapper.Cmd, func(arg string, _ int, _ []string) string {
				return expander.Expand(arg)
			}),
			Env:  map[string]string{},
			With: map[string]string{},
		}
		id = js.Coalesce(expandedStep.Id, expandedStep.Name, expandedStep.Uses)
	)

	for k, v := range e.StepWrapper.Env {
		expandedStep.Env[k] = expander.Expand(v)
	}

	for k, v := range e.StepWrapper.With {
		expandedStep.With[k] = expander.Expand(v)
	}

	e.GlobalContext.InputsContext = expandedStep.With
	e.GlobalContext.AddEnv(expandedStep.Env)
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
				fmt.Sprintf("%s=%s", actions.EnvVarEnv, "TODO"),
				fmt.Sprintf("%s=%s", actions.EnvVarPath, "TODO"),
			},
			e.GlobalContext.EnvArr()...,
		),
		Mounts: append(runtime.NetworkMounts, e.StepWrapper.ExtraMounts...),
	}

	for k, v := range e.StepWrapper.State {
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
