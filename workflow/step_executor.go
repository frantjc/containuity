package workflow

import (
	"os"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/conf"
)

func NewStepExecutor(s *Step, opts ...ExecOpt) (Executor, error) {
	ex := &jobExecutor{
		stdout: os.Stdout,
		stderr: os.Stderr,

		globalContext: actions.EmptyContext(),

		runnerImage: conf.DefaultRunnerImage,

		steps: []*Step{s},

		ctxOpts: []actions.CtxOpt{
			actions.WithEnv(s.Env),
			func(gc *actions.GlobalContext) error {
				if gc.GitHubContext == nil {
					gc.GitHubContext = &actions.GitHubContext{}
				}
				if gc.RunnerContext == nil {
					gc.RunnerContext = &actions.RunnerContext{}
				}
				gc.GitHubContext.ActionPath = containerActionPath
				gc.GitHubContext.Workspace = containerWorkspace
				gc.RunnerContext.Temp = containerRunnerTemp
				gc.RunnerContext.ToolCache = containerRunnerToolCache
				return nil
			},
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
