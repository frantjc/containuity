package workflow

import (
	"bytes"
	"io"

	"github.com/frantjc/sequence/meta"
)

type runOpts struct {
	path        string
	jobName     string
	job         *Job
	workflow    *Workflow
	githubToken string
	stdout      io.Writer
	stderr      io.Writer
	verbose     bool
	image       string
}

type RunOpt func(*runOpts) error

func WithJobName(j string) RunOpt {
	return func(ro *runOpts) error {
		ro.jobName = j
		return nil
	}
}

func WithJob(j *Job) RunOpt {
	return func(ro *runOpts) error {
		ro.job = j
		return nil
	}
}

func WithWorkflow(w *Workflow) RunOpt {
	return func(ro *runOpts) error {
		ro.workflow = w
		return nil
	}
}

func WithGitHubToken(token string) RunOpt {
	return func(ro *runOpts) error {
		ro.githubToken = token
		return nil
	}
}

func WithStdout(stdout io.Writer) RunOpt {
	return func(ro *runOpts) error {
		ro.stdout = stdout
		return nil
	}
}

func WithStderr(stderr io.Writer) RunOpt {
	return func(ro *runOpts) error {
		ro.stderr = stderr
		return nil
	}
}

func WithVerbose(ro *runOpts) error {
	ro.verbose = true
	return nil
}

func WithImage(image string) RunOpt {
	return func(ro *runOpts) error {
		ro.image = image
		return nil
	}
}

func newRunOpts(opts ...RunOpt) (*runOpts, error) {
	var (
		stdout = new(bytes.Buffer)
		ro     = &runOpts{
			workflow: &Workflow{},
			job:      &Job{},
			path:     ".",
			stdout:   stdout,
			stderr:   stdout,
			verbose:  false,
		}
	)
	for _, opt := range opts {
		err := opt(ro)
		if err != nil {
			return nil, err
		}
	}

	if jobImage, ok := ro.job.Container.(string); ok {
		ro.image = jobImage
	}

	if ro.image == "" {
		ro.image = meta.Image()
	}

	if ro.job != nil {
		ro.jobName = ro.job.Name
	}

	return ro, nil
}
