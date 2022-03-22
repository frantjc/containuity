package actions

import (
	"context"

	"github.com/go-git/go-git/v5"
)

type Vars struct {
	Env            *Env
	ActionsContext *ActionsContext
}

func NewVarsFromPath(ctx context.Context, path string, opts ...VarsOpt) (*Vars, error) {
	vopts := defaultVarsOpts()
	for _, opt := range opts {
		err := opt(vopts)
		if err != nil {
			return nil, err
		}
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	env, err := newEnvFromRepository(ctx, repo, vopts)
	if err != nil {
		return nil, err
	}

	actx, err := newCtxFromRepository(ctx, repo, vopts)
	if err != nil {
		return nil, err
	}

	return &Vars{
		Env:            env,
		ActionsContext: actx,
	}, nil
}
