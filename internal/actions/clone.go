package actions

import (
	"context"
	"net/url"
	"path/filepath"

	"github.com/frantjc/sequence/internal/github"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type cloneOpts struct {
	path      string
	insecure  bool
	gitHubURL *url.URL
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

func defaultCloneOps() *cloneOpts {
	return &cloneOpts{
		path:      ".",
		gitHubURL: github.DefaultURL,
	}
}

func Clone(u *Uses, opts ...CloneOpt) (*Action, error) {
	return CloneContext(context.Background(), u, opts...)
}

func CloneContext(ctx context.Context, u *Uses, opts ...CloneOpt) (*Action, error) {
	copts := defaultCloneOps()
	for _, opt := range opts {
		err := opt(copts)
		if err != nil {
			return nil, err
		}
	}

	cloneURL := copts.gitHubURL
	cloneURL.Path = u.Path

	repo, err := git.PlainCloneContext(ctx, copts.path, false, &git.CloneOptions{
		URL:               cloneURL.String(),
		ReferenceName:     plumbing.NewTagReferenceName(u.Version),
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		InsecureSkipTLS:   copts.insecure,
		Tags:              git.AllTags,
	})
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	com, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	var f *object.File
	f, err = com.File(filepath.Join(u.Path, "action.yml"))
	if err == object.ErrFileNotFound {
		f, err = com.File(filepath.Join(u.Path, "action.yaml"))
		if err != nil {
			return nil, ErrNotAnAction
		}
	} else if err != nil {
		return nil, err
	}

	r, err := f.Reader()
	if err != nil {
		return nil, err
	}

	return NewActionFromReader(r)
}
