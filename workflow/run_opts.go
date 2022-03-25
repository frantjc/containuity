package workflow

import (
	"bytes"
	"io"

	"github.com/frantjc/sequence/conf"
)

type runOpts struct {
	repository  string
	jobName     string
	job         *Job
	workflow    *Workflow
	githubToken string
	stdout      io.Writer
	stderr      io.Writer
	verbose     bool
	actionImage string
	runnerImage string
	workdir     string
	secrets     map[string]string
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

		if j.Name != "" {
			ro.jobName = j.Name
		}

		if jobImage, ok := ro.job.Container.(string); ok {
			ro.runnerImage = jobImage
		}

		return nil
	}
}

func WithWorkflow(w *Workflow) RunOpt {
	return func(ro *runOpts) error {
		ro.workflow = w
		if ro.job != nil && ro.jobName == "" {
			for name, job := range ro.workflow.Jobs {
				if ro.job == &job {
					ro.jobName = name
				}
			}
		}
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

func WithActionImage(image string) RunOpt {
	return func(ro *runOpts) error {
		ro.actionImage = image
		return nil
	}
}

func WithRunnerImage(image string) RunOpt {
	return func(ro *runOpts) error {
		ro.runnerImage = image
		return nil
	}
}

func WithWorkdir(workdir string) RunOpt {
	return func(ro *runOpts) error {
		ro.workdir = workdir
		return nil
	}
}

func WithRepository(repository string) RunOpt {
	return func(ro *runOpts) error {
		ro.repository = repository
		return nil
	}
}

func WithSecrets(secrets map[string]string) RunOpt {
	return func(ro *runOpts) error {
		if ro.secrets == nil {
			ro.secrets = secrets
		} else {
			for k, v := range secrets {
				ro.secrets[k] = v
			}
		}
		return nil
	}
}

func newRunOpts(opts ...RunOpt) (*runOpts, error) {
	var (
		buf = new(bytes.Buffer)
		ro  = &runOpts{
			workflow:    &Workflow{},
			job:         &Job{},
			repository:  ".",
			stdout:      buf,
			stderr:      buf,
			workdir:     ".",
			runnerImage: conf.DefaultRunnerImage,
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
