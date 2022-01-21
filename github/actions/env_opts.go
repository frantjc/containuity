package actions

import "os/user"

type envOpts struct {
	remote       string
	branch       string
	workdir      string
	workflow     string
	runID        int
	runNumber    int
	job          string
	refProtected bool
	headRef      string
	baseRef      string
	runnerName   string
	token        string
}

type EnvOpt func(e *envOpts) error

func defaultEnvOpts() *envOpts {
	u, _ := user.Current()
	return &envOpts{
		remote:     defaultRemote,
		branch:     defaultBranch,
		workdir:    "/sqnc",
		runnerName: u.Username,
	}
}
