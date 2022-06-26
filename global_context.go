package sequence

import "github.com/frantjc/sequence/github/actions"

func defaultGlobalContextOpts() []actions.CtxOpt {
	return []actions.CtxOpt{
		func(gc *actions.GlobalContext) error {
			gc.GitHubContext.ActionPath = actionPath
			gc.GitHubContext.Workspace = workspace
			gc.RunnerContext.Temp = runnerTemp
			gc.RunnerContext.ToolCache = runnerToolCache
			return nil
		},
	}
}
