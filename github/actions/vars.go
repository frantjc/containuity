package actions

import "github.com/go-git/go-git/v5"

type Vars struct {
	Env            *Env
	ActionsContext *ActionsContext
}

func NewVarsFromPath(path string, opts ...VarsOpt) (*Vars, error) {
	vopts := defaultVarsOpts()
	for _, opt := range opts {
		err := opt(vopts)
		if err != nil {
			return nil, err
		}
	}

	eopts := defaultEnvOpts()
	for _, opt := range vopts.eopts {
		err := opt(eopts)
		if err != nil {
			return nil, err
		}
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	env, err := newEnvFromRepository(repo, eopts)
	if err != nil {
		return nil, err
	}

	copts := defaultCtxOpts()
	for _, opt := range vopts.copts {
		err := opt(copts)
		if err != nil {
			return nil, err
		}
	}

	ctx, err := newCtxFromRepository(repo, copts)
	if err != nil {
		return nil, err
	}

	return &Vars{
		Env:            env,
		ActionsContext: ctx,
	}, nil
}
