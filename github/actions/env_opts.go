package actions

import "os"

type envOpts struct {
	remote    string
	branch    string
	tmpdir    string
	workdir   string
	toolcache string
	runID     int
	runNumber int
}

type EnvOpt func(e *envOpts) error

func defaultEnvOpts() *envOpts {
	cachedir, _ := os.UserCacheDir()

	return &envOpts{
		remote:    "origin",
		branch:    "main",
		tmpdir:    os.TempDir(),
		toolcache: cachedir,
	}
}
