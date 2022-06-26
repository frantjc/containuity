package sequence

import (
	"io"
	"os"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/runtime"
)

type ExecutorOpt func(*Executor) error

func WithVerbose(e *Executor) error {
	e.Verbose = true
	return nil
}

func WithStreams(stdin io.Reader, stdout, stderr io.Writer) ExecutorOpt {
	return func(e *Executor) error {
		e.Stdin = stdin
		e.Stdout = stdout
		e.Stderr = stderr
		return nil
	}
}

var WithStdio = WithStreams(os.Stdin, os.Stdout, os.Stderr)

func WithRuntime(runtime runtime.Runtime) ExecutorOpt {
	return func(e *Executor) error {
		e.Runtime = runtime
		return nil
	}
}

func WithRunnerImage(image runtime.Image) ExecutorOpt {
	return func(e *Executor) error {
		e.RunnerImage = image
		return nil
	}
}

func WithGlobalContext(gc *actions.GlobalContext) ExecutorOpt {
	return func(e *Executor) error {
		e.GlobalContext = gc
		for _, opt := range defaultGlobalContextOpts() {
			if err := opt(e.GlobalContext); err != nil {
				return err
			}
		}
		return nil
	}
}

func OnImagePull(hooks ...Hook[runtime.Image]) ExecutorOpt {
	return func(e *Executor) error {
		e.OnImagePull = append(e.OnImagePull, hooks...)
		return nil
	}
}

func OnContainerCreate(hooks ...Hook[runtime.Container]) ExecutorOpt {
	return func(e *Executor) error {
		e.OnContainerCreate = append(e.OnContainerCreate, hooks...)
		return nil
	}
}

func OnVolumeCreate(hooks ...Hook[runtime.Volume]) ExecutorOpt {
	return func(e *Executor) error {
		e.OnVolumeCreate = append(e.OnVolumeCreate, hooks...)
		return nil
	}
}

func OnWorkflowCommand(hooks ...Hook[*actions.WorkflowCommand]) ExecutorOpt {
	return func(e *Executor) error {
		e.OnWorkflowCommand = append(e.OnWorkflowCommand, hooks...)
		return nil
	}
}
