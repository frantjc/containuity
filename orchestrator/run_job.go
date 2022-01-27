package orchestrator

import (
	"context"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
)

func RunJob(ctx context.Context, r runtime.Runtime, j *sequence.Job, opts ...OrchOpt) error {
	for _, step := range j.Steps {
		err := RunStep(ctx, r, &step, append(opts, WithJob(j))...)
		if err != nil {
			return err
		}
	}

	return nil
}
