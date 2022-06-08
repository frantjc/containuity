package sequence

import (
	"os"

	workflowv1 "github.com/frantjc/sequence/workflow/v1"
)

type runOpts struct {
	job         *workflowv1.Job
	jobName     string
	workflow    *workflowv1.Workflow
	runnerImage string
	verbose     bool
	repository  string
}

func defaultRunOpts() *runOpts {
	wd, _ := os.Getwd()
	return &runOpts{
		repository: wd,
	}
}

type RunOpt func(*runOpts) error

func WithJob(j *workflowv1.Job) RunOpt {
	return func(ro *runOpts) error {
		ro.job = j

		ro.runnerImage = j.Container.GetImage()

		if j.Name != "" {
			ro.jobName = j.Name
		}

		return nil
	}
}

func WithWorkflow(w *workflowv1.Workflow) RunOpt {
	return func(ro *runOpts) error {
		ro.workflow = w
		return nil
	}
}

func WithVerbose(ro *runOpts) error {
	ro.verbose = true
	return nil
}

func WithRepository(r string) RunOpt {
	return func(ro *runOpts) error {
		ro.repository = r
		return nil
	}
}

func WithRunnerImage(i string) RunOpt {
	return func(ro *runOpts) error {
		ro.runnerImage = i
		return nil
	}
}
