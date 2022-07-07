package actions

import (
	"net/url"
	"path/filepath"

	"github.com/frantjc/sequence/pkg/github"
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
	return func(co *cloneOpts) error {
		co.path = filepath.Clean(p)
		return nil
	}
}

func WithGitHubURL(u string) CloneOpt {
	return func(co *cloneOpts) error {
		var err error
		co.githubURL, err = url.Parse(u)
		return err
	}
}

func WithInsecure(co *cloneOpts) error {
	co.insecure = true
	return nil
}
