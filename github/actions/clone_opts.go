package actions

import (
	"net/url"
	"path/filepath"

	"github.com/frantjc/sequence/github"
)

type cloneOpts struct {
	path      string
	insecure  bool
	githubURL *url.URL
}

func defaultCloneOpts() *cloneOpts {
	return &cloneOpts{
		path:      ".",
		githubURL: github.DefaultURL,
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
		copts.githubURL, err = url.Parse(u)
		return err
	}
}

func WithInsecure() CloneOpt {
	return func(copts *cloneOpts) error {
		copts.insecure = true
		return nil
	}
}
