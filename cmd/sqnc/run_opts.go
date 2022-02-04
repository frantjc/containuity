package main

import (
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
)

type runOpts struct {
	path        string
	jobName     string
	job         *sequence.Job
	workflow    *sequence.Workflow
	sopts       []runtime.SpecOpt
	gitHubToken string
}

type runOpt func(*runOpts) error

func withJobName(j string) runOpt {
	return func(ro *runOpts) error {
		ro.jobName = j
		return nil
	}
}

func withJob(j *sequence.Job) runOpt {
	return func(ro *runOpts) error {
		ro.job = j
		return nil
	}
}

func withWorkflow(w *sequence.Workflow) runOpt {
	return func(ro *runOpts) error {
		ro.workflow = w
		return nil
	}
}

func withGitHubToken(token string) runOpt {
	return func(ro *runOpts) error {
		ro.gitHubToken = token
		return nil
	}
}
