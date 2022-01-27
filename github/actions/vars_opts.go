package actions

import "os/user"

type varsOpts struct {
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

type VarsOpt func(v *varsOpts) error

func defaultVarsOpts() *varsOpts {
	u, _ := user.Current()
	return &varsOpts{
		remote:     defaultRemote,
		branch:     defaultBranch,
		workdir:    "/sqnc",
		runnerName: u.Username,
	}
}

func WithToken(token string) VarsOpt {
	return func(vo *varsOpts) error {
		vo.token = token
		return nil
	}
}
