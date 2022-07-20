package sequence

import (
	"context"
	"fmt"
	"path"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
)

type stepWrapperExecutor struct {
	*executor
	stepWrapper        *stepWrapper
	stopCommandsTokens map[string]bool
}

// WorkflowCommandWriterCallback swallows the bytes of
// _all_ non-stopped workflow commands.
// "Non-stopped" means not between a "stop-commands"
// workflow command and its end token.
// However, it handles the functionality of all workflow commands.
// It is up to the caller to handle the logging of workflow
// commands if they so choose.
func (e *stepWrapperExecutor) WorkflowCommandWriterCallback(wc *actions.WorkflowCommand) []byte {
	event := &Event[*actions.WorkflowCommand]{
		Type:          wc,
		GlobalContext: e.GlobalContext,
	}

	if _, ok := e.stopCommandsTokens[wc.Command]; ok {
		e.OnWorkflowCommand.Invoke(event)
		e.stopCommandsTokens[wc.Command] = false
		return make([]byte, 0)
	}

	for _, stop := range e.stopCommandsTokens {
		if stop {
			return []byte(wc.String())
		}
	}

	e.OnWorkflowCommand.Invoke(event)

	switch wc.Command {
	case actions.CommandSetOutput:
		if e.GlobalContext.StepsContext[e.stepWrapper.step.GetId()] == nil {
			e.GlobalContext.StepsContext[e.stepWrapper.step.GetId()] = &actions.StepsContext{
				Outputs: map[string]string{},
			}
		}

		e.GlobalContext.StepsContext[e.stepWrapper.step.GetId()].Outputs[wc.GetName()] = wc.Value
	case actions.CommandStopCommands:
		e.stopCommandsTokens[wc.Value] = true
	case actions.CommandSaveState:
		e.stepWrapper.state[wc.GetName()] = wc.Value
	}

	return make([]byte, 0)
}

func (e *stepWrapperExecutor) Execute(ctx context.Context) error {
	var (
		expander     = actions.NewExpander(e.GlobalContext.GetString)
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
			// expand these later as it can't be easily reduced into a one-liner
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
		Entrypoint: []string{paths.Shim, "-e"},
		Cwd:        e.GlobalContext.GitHubContext.Workspace,
		Env: append(
			[]string{
				"SQNC=true",
				"SEQUENCE=true",
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
			e.StreamIn,
			actions.NewWorkflowCommandWriter(e.WorkflowCommandWriterCallback, e.StreamOut),
			actions.NewWorkflowCommandWriter(e.WorkflowCommandWriterCallback, e.StreamErr),
		),
	)
}
