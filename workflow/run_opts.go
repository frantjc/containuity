package workflow

import (
	"bytes"
	"io"

	"github.com/frantjc/sequence/conf"
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
	workdir     string
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

		if ro.jobName == "" {
			ro.jobName = j.Name
		}

		if jobImage, ok := ro.job.Container.(string); ok {
			ro.image = jobImage
		} else if container, ok := ro.job.Container.(*Container); ok {
			ro.image = container.Image
		}

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

func WithRunnerImage(image string) RunOpt {
	return func(ro *runOpts) error {
		ro.image = image
		return nil
	}
}

func WithWorkdir(workdir string) RunOpt {
	return func(ro *runOpts) error {
		ro.workdir = workdir
		return nil
	}
}

func newRunOpts(opts ...RunOpt) (*runOpts, error) {
	var (
		buf = new(bytes.Buffer)
		ro  = &runOpts{
			workflow: &Workflow{},
			job:      &Job{},
			path:     ".",
			stdout:   buf,
			stderr:   buf,
			verbose:  false,
			workdir:  ".",
			image:    conf.DefaultRunnerImage,
		}
	)
	for _, opt := range opts {
		err := opt(ro)
		if err != nil {
			return nil, err
		}
	}

	return ro, nil
}
