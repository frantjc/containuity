package sequence

import (
	"os"

	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/workflow"
)

type runOpts struct {
	job         *workflow.Job
	jobName     string
	workflow    *workflow.Workflow
	runnerImage string
	verbose     bool
	repository  string
}

func defaultRunOpts() *runOpts {
	wd, _ := os.Getwd()
	return &runOpts{
		repository:  wd,
		runnerImage: conf.DefaultRunnerImage,
	}
}

type RunOpt func(*runOpts) error

func WithJob(j *workflow.Job) RunOpt {
	return func(ro *runOpts) error {
		ro.job = j

		if j.Name != "" {
			ro.jobName = j.Name
		}

		if jobContainer, ok := j.Container.(*workflow.Container); ok {
			ro.runnerImage = jobContainer.Image
		} else if jobImage, ok := j.Container.(string); ok {
			ro.runnerImage = jobImage
		}

		return nil
	}
}

func WithWorkflow(w *workflow.Workflow) RunOpt {
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
