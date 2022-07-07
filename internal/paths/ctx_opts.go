package paths

import "github.com/frantjc/sequence/pkg/github/actions"

func GlobalContextOpts() []actions.CtxOpt {
	return []actions.CtxOpt{
		func(gc *actions.GlobalContext) error {
			gc.GitHubContext.ActionPath = Action
			gc.GitHubContext.Workspace = Workspace
			gc.RunnerContext.Temp = RunnerTemp
			gc.RunnerContext.ToolCache = RunnerToolCache
			return nil
		},
	}
}
