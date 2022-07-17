package sequence

import (
	"io"
	"os"

	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
)

type ExecutorOpt func(*executor) error

func WithID(id string) ExecutorOpt {
	return func(e *executor) error {
		e.ID = id
		return nil
	}
}

func WithVerbose(e *executor) error {
	e.Verbose = true
	return nil
}

func WithStreams(stdin io.Reader, stdout, stderr io.Writer) ExecutorOpt {
	return func(e *executor) error {
		e.Stdin = stdin
		e.Stdout = stdout
		e.Stderr = stderr
		return nil
	}
}

var WithStdio = WithStreams(os.Stdin, os.Stdout, os.Stderr)

func WithRuntime(runtime runtime.Runtime) ExecutorOpt {
	return func(e *executor) error {
		e.Runtime = runtime
		return nil
	}
}

func WithRunnerImage(image runtime.Image) ExecutorOpt {
	return func(e *executor) error {
		e.RunnerImage = image
		return nil
	}
}

func WithGlobalContext(gc *actions.GlobalContext) ExecutorOpt {
	return func(e *executor) error {
		e.GlobalContext = gc
		for _, opt := range paths.GlobalContextOpts() {
			if err := opt(e.GlobalContext); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithJobName(jobName string) ExecutorOpt {
	return func(e *executor) error {
		return actions.WithJobName(jobName)(e.GlobalContext)
	}
}

func WithWorkflowName(workflowName string) ExecutorOpt {
	return func(e *executor) error {
		return actions.WithWorkflowName(workflowName)(e.GlobalContext)
	}
}

func OnImagePull(hooks ...Hook[runtime.Image]) ExecutorOpt {
	return func(e *executor) error {
		e.OnImagePull = append(e.OnImagePull, hooks...)
		return nil
	}
}

func OnContainerCreate(hooks ...Hook[runtime.Container]) ExecutorOpt {
	return func(e *executor) error {
		e.OnContainerCreate = append(e.OnContainerCreate, hooks...)
		return nil
	}
}

func OnVolumeCreate(hooks ...Hook[runtime.Volume]) ExecutorOpt {
	return func(e *executor) error {
		e.OnVolumeCreate = append(e.OnVolumeCreate, hooks...)
		return nil
	}
}

func OnWorkflowCommand(hooks ...Hook[*actions.WorkflowCommand]) ExecutorOpt {
	return func(e *executor) error {
		e.OnWorkflowCommand = append(e.OnWorkflowCommand, hooks...)
		return nil
	}
}

func OnStepStart(hooks ...Hook[*Step]) ExecutorOpt {
	return func(e *executor) error {
		e.OnStepStart = append(e.OnStepStart, hooks...)
		return nil
	}
}

func OnStepFinish(hooks ...Hook[*Step]) ExecutorOpt {
	return func(e *executor) error {
		e.OnStepFinish = append(e.OnStepFinish, hooks...)
		return nil
	}
}

func OnJobStart(hooks ...Hook[*Job]) ExecutorOpt {
	return func(e *executor) error {
		e.OnJobStart = append(e.OnJobStart, hooks...)
		return nil
	}
}

func OnJobFinish(hooks ...Hook[*Job]) ExecutorOpt {
	return func(e *executor) error {
		e.OnJobFinish = append(e.OnJobFinish, hooks...)
		return nil
	}
}

func OnWorkflowStart(hooks ...Hook[*Workflow]) ExecutorOpt {
	return func(e *executor) error {
		e.OnWorkflowStart = append(e.OnWorkflowStart, hooks...)
		return nil
	}
}

func OnWorkflowFinish(hooks ...Hook[*Workflow]) ExecutorOpt {
	return func(e *executor) error {
		e.OnWorkflowFinish = append(e.OnWorkflowFinish, hooks...)
		return nil
	}
}
