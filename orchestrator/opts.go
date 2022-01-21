package orchestrator

import (
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
)

type orchOpts struct {
	path     string
	jobName  string
	job      *sequence.Job
	workflow *sequence.Workflow
	sopts []runtime.SpecOpt
}

type OrchOpt func(*orchOpts) error

func WithSpecOpts(opts ...runtime.SpecOpt) OrchOpt {
	return func(oo *orchOpts) error {
		oo.sopts = opts
		return nil
	}
}

func WithPath(path string) OrchOpt {
	return func(oo *orchOpts) error {
		oo.path = path
		return nil
	}
}

func WithJobName(j string) OrchOpt {
	return func(oo *orchOpts) error {
		oo.jobName = j
		return nil
	}
}

func WithJob(j *sequence.Job) OrchOpt {
	return func(oo *orchOpts) error {
		oo.job = j
		return nil
	}
}

func WithWorkflow(w *sequence.Workflow) OrchOpt {
	return func(oo *orchOpts) error {
		oo.workflow = w
		return nil
	}
}
