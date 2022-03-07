package workflow

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

func RunJob(ctx context.Context, r runtime.Runtime, j *Job, opts ...RunOpt) error {
	ro, err := newRunOpts(append(opts, WithJob(j))...)
	if err != nil {
		return err
	}

	for _, step := range j.Steps {
		err := runStep(ctx, r, &step, ro)
		if err != nil {
			return err
		}
	}

	return nil
}
