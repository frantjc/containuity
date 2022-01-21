package actions

type ctxOpts struct {
	remote string
	branch string
}

type CtxOpt func(e *ctxOpts) error

func defaultCtxOpts() *ctxOpts {
	return &ctxOpts{
		remote: defaultRemote,
		branch: defaultBranch,
	}
}
