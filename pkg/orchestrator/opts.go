package orchestrator

import (
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/container"
)

type runOpts struct {
	workflow *sequence.Workflow
	job      *sequence.Job
	path     string
	mounts   []container.Mount
	env      []string
}

type RunOpt func(*runOpts) error

func WithWorkflow(workflow *sequence.Workflow) RunOpt {
	return func(ro *runOpts) error {
		ro.workflow = workflow
		return nil
	}
}

func WithJob(job *sequence.Job) RunOpt {
	return func(ro *runOpts) error {
		ro.job = job
		return nil
	}
}

func WithPath(path string) RunOpt {
	return func(ro *runOpts) error {
		ro.path = path
		return nil
	}
}

func WithMounts(mounts []container.Mount) RunOpt {
	return func(ro *runOpts) error {
		ro.mounts = mounts
		return nil
	}
}

func WithEnv(env []string) RunOpt {
	return func(ro *runOpts) error {
		ro.env = env
		return nil
	}
}
