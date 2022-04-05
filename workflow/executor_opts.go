package workflow

import (
	"context"
	"io"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/runtime"
)

type ExecOpt func(*jobExecutor) error

func WithGlobalContext(globalContext *actions.GlobalContext) ExecOpt {
	return func(e *jobExecutor) error {
		for _, opt := range e.ctxOpts {
			if err := opt(globalContext); err != nil {
				return err
			}
		}

		e.globalContext = globalContext

		return nil
	}
}

func WithRepository(repository string) ExecOpt {
	return func(e *jobExecutor) (err error) {
		e.globalContext, err = actions.NewContextFromPath(context.Background(), repository, e.ctxOpts...)
		return
	}
}

func WithStdout(stdout io.Writer) ExecOpt {
	return func(e *jobExecutor) error {
		e.stdout = stdout
		return nil
	}
}

func WithStderr(stderr io.Writer) ExecOpt {
	return func(e *jobExecutor) error {
		e.stderr = stderr
		return nil
	}
}

func WithRunnerImage(runnerImage string) ExecOpt {
	return func(e *jobExecutor) error {
		e.runnerImage = runnerImage
		return nil
	}
}

func WithGitHubToken(token string) ExecOpt {
	return func(e *jobExecutor) error {
		e.globalContext.GitHubContext.Token = token
		e.ctxOpts = append(e.ctxOpts, actions.WithToken(token))
		return nil
	}
}

func WithSecrets(secrets map[string]string) ExecOpt {
	return func(e *jobExecutor) error {
		for k, v := range secrets {
			e.globalContext.SecretsContext[k] = v
		}
		e.ctxOpts = append(e.ctxOpts, actions.WithSecrets(secrets))
		return nil
	}
}

func WithJob(j *Job) ExecOpt {
	return func(e *jobExecutor) error {
		if jobImage, ok := j.Container.(string); ok {
			e.runnerImage = jobImage
		}

		if j.Name != "" {
			e.globalContext.GitHubContext.Job = j.Name
			e.ctxOpts = append(e.ctxOpts, actions.WithJobName(j.Name))
		}

		return nil
	}
}

func WithWorkflow(w *Workflow) ExecOpt {
	return func(e *jobExecutor) error {
		if w.Name != "" {
			e.globalContext.GitHubContext.Workflow = w.Name
			e.ctxOpts = append(e.ctxOpts, actions.WithWorkflowName(w.Name))
		}

		return nil
	}
}

func WithVerbose(e *jobExecutor) error {
	e.verbose = true
	return nil
}

func WithRuntime(r runtime.Runtime) ExecOpt {
	return func(e *jobExecutor) error {
		e.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) ExecOpt {
	return func(e *jobExecutor) (err error) {
		e.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime = WithRuntimeName()

func WithWorkdir(workdir string) ExecOpt {
	return func(e *jobExecutor) error {
		e.workdir = workdir
		return nil
	}
}
