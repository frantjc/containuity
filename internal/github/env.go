package github

import (
	"github.com/go-git/go-git/v5"
)

type envOpts struct {
	remote string
	branch string
}

type EnvOpt func(e *envOpts) error

func defaultEnvOpts() *envOpts {
	return &envOpts{
		remote: "origin",
		branch: "main",
	}
}

// WIP
func EnvFromRepository(path string, opts ...EnvOpt) (map[string]string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	eopts := defaultEnvOpts()
	if ref.Name().IsBranch() {
		eopts.branch = ref.String()
	}

	for _, opt := range opts {
		err := opt(eopts)
		if err != nil {
			return nil, err
		}
	}

	_, err = repo.Branch(eopts.branch)
	if err != nil {
		return nil, err
	}
	
	_, err = repo.Remote(eopts.remote)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"CI": "true",
	}, nil
}
