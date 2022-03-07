package workflow

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

func RunWorkflow(ctx context.Context, r runtime.Runtime, w *Workflow, opts ...RunOpt) error {
	for name, job := range w.Jobs {
		jobName := name
		if job.Name != "" {
			jobName = job.Name
		}

		err := RunJob(ctx, r, &job, append(opts, WithJobName(jobName), WithJob(&job), WithWorkflow(w))...)
		if err != nil {
			return err
		}
	}

	return nil
}
