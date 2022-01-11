package actions

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/frantjc/sequence/github"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
)

type cloneOpts struct {
	path      string
	insecure  bool
	gitHubURL *url.URL
}

func defaultCloneOps() *cloneOpts {
	return &cloneOpts{
		path:      ".",
		gitHubURL: github.DefaultURL,
	}
}

type CloneOpt func(*cloneOpts) error

func WithPath(p string) CloneOpt {
	return func(copts *cloneOpts) error {
		copts.path = filepath.Clean(p)
		return nil
	}
}

func WithGitHubURL(u string) CloneOpt {
	return func(copts *cloneOpts) error {
		var err error
		copts.gitHubURL, err = url.Parse(u)
		return err
	}
}

func WithInsecure() CloneOpt {
	return func(copts *cloneOpts) error {
		copts.insecure = true
		return nil
	}
}

func Clone(r Reference, opts ...CloneOpt) (*Action, error) {
	return CloneContext(context.Background(), r, opts...)
}

func CloneContext(ctx context.Context, r Reference, opts ...CloneOpt) (*Action, error) {
	var (
		copts  = defaultCloneOps()
		cloned = false
	)
	for _, opt := range opts {
		err := opt(copts)
		if err != nil {
			return nil, err
		}
	}

	cloneURL := copts.gitHubURL
	cloneURL.Path = fullRepository(r)
	clopts := &git.CloneOptions{
		URL:               cloneURL.String(),
		ReferenceName:     plumbing.NewTagReferenceName(r.Version()),
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		InsecureSkipTLS:   copts.insecure,
		Tags:              git.TagFollowing,
	}
	log.Debug().Msgf("checking %s for existing %s %s repository", copts.path, clopts.URL, clopts.ReferenceName)
	// check if the desired action is already cloned to the given path
	repo, err := git.PlainOpen(copts.path)
	if err == nil {
		remotes, err := repo.Remotes()
		if err != nil {
			return nil, err
		}

		// check if the repository cloned to the given path is the repository of the desired action
		for _, remote := range remotes {
			for _, confURL := range remote.Config().URLs {
				if confURL == clopts.URL {
					log.Debug().Msgf("%s already cloned to %s", clopts.URL, copts.path)
					cloned = true
				}
			}
		}

		// check if the repository HEAD at the given path is the desired version of the action
		if cloned {
			ref, err := repo.Head()
			if err != nil {
				cloned = ref.String() == r.Version()
			}
		}
	}

	if !cloned {
		log.Debug().Msgf("%s %s was not cloned to %s, cleaning %s", clopts.URL, clopts.ReferenceName, copts.path, copts.path)
		err = os.RemoveAll(copts.path)
		if err != nil {
			return nil, err
		}

		log.Debug().Msgf("cloning %s %s to %s", clopts.URL, clopts.ReferenceName, copts.path)
		repo, err = git.PlainCloneContext(ctx, copts.path, false, clopts)
		if err != nil {
			log.Debug().Msgf("cloning %s with ref assumed as tag, falling back to branch", cloneURL.String())
			clopts.ReferenceName = plumbing.NewBranchReferenceName(r.Version())
			log.Debug().Msgf("cloning %s %s to %s", clopts.URL, clopts.ReferenceName, copts.path)
			repo, err = git.PlainCloneContext(ctx, copts.path, false, clopts)
			if err != nil {
				return nil, err
			}
		}
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("searching %s HEAD for action.yml or action.yaml", copts.path)
	var f *object.File
	f, err = commit.File(filepath.Join(r.Path(), "action.yml"))
	if errors.Is(err, object.ErrFileNotFound) {
		f, err = commit.File(filepath.Join(r.Path(), "action.yaml"))
		if err != nil {
			return nil, ErrNotAnAction
		}
	} else if err != nil {
		return nil, err
	}

	log.Debug().Msgf("found action from %s %s", clopts.URL, clopts.ReferenceName)
	a, err := f.Reader()
	if err != nil {
		return nil, err
	}

	return NewActionFromReader(a)
}
