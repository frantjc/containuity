package plan

import (
	"github.com/frantjc/sequence"
)

type planOpts struct {
	path     string
	jobName  string
	job      *sequence.Job
	workflow *sequence.Workflow
}

type PlanOpt func(*planOpts) error

func WithPath(path string) PlanOpt {
	return func(po *planOpts) error {
		po.path = path
		return nil
	}
}

func WithJobName(j string) PlanOpt {
	return func(po *planOpts) error {
		po.jobName = j
		return nil
	}
}

func WithJob(j *sequence.Job) PlanOpt {
	return func(po *planOpts) error {
		po.job = j
		return nil
	}
}

func WithWorkflow(w *sequence.Workflow) PlanOpt {
	return func(po *planOpts) error {
		po.workflow = w
		return nil
	}
}
