package actions

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/frantjc/sequence/pkg/github/actions/uses"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Clone(u *uses.Uses, opts ...CloneOpt) (*Metadata, error) {
	return CloneContext(context.Background(), u, opts...)
}

func CloneContext(ctx context.Context, u *uses.Uses, opts ...CloneOpt) (*Metadata, error) {
	copts := defaultCloneOpts()
	for _, opt := range opts {
		err := opt(copts)
		if err != nil {
			return nil, err
		}
	}

	cloneURL := copts.githubURL
	cloneURL.Path = u.FullRepository()
	clopts := &git.CloneOptions{
		URL:               cloneURL.String(),
		ReferenceName:     plumbing.NewTagReferenceName(u.Version),
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		InsecureSkipTLS:   copts.insecure,
		Tags:              git.TagFollowing,
	}
	repo, err := git.PlainCloneContext(ctx, copts.path, false, clopts)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		repo, err = git.PlainOpen(copts.path)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		clopts.ReferenceName = plumbing.NewBranchReferenceName(u.Version)
		repo, err = git.PlainCloneContext(ctx, copts.path, false, clopts)
		if err != nil {
			return nil, err
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

	var f *object.File
	f, err = commit.File(filepath.Join(u.Path, "action.yml"))
	if errors.Is(err, object.ErrFileNotFound) {
		f, err = commit.File(filepath.Join(u.Path, "action.yaml"))
		if err != nil {
			return nil, ErrNotAnAction
		}
	} else if err != nil {
		return nil, err
	}

	m, err := f.Reader()
	if err != nil {
		return nil, err
	}

	return NewMetadataFromReader(m)
}
