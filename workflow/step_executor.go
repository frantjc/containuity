package workflow

import (
	"os"

	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/github/actions"
)

func NewStepExecutor(s *Step, opts ...ExecOpt) (Executor, error) {
	ex := &jobExecutor{
		stdout: os.Stdout,
		stderr: os.Stderr,

		globalContext: actions.EmptyContext(),

		runnerImage: conf.DefaultRunnerImage,

		steps: []*Step{s},

		ctxOpts: []actions.CtxOpt{
			actions.WithWorkdir(containerWorkdir),
			actions.WithEnv(s.Env),
		},

		states: map[string]map[string]string{},
	}

	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}

	return ex, nil
}
