package orchestrator

import (
	"context"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
)

func RunWorkflow(ctx context.Context, r runtime.Runtime, w *sequence.Workflow, opts ...OrchOpt) error {
	for name, job := range w.Jobs {
		err := RunJob(ctx, r, &job, append(opts, WithJobName(name), WithJob(&job))...)
		if err != nil {
			return err
		}
	}

	return nil
}
