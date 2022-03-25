package workflow

import (
	"context"
	"fmt"

	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
)

func RunJob(ctx context.Context, r runtime.Runtime, j *Job, opts ...RunOpt) (context.Context, error) {
	ro, err := newRunOpts(append(opts, WithJob(j))...)
	if err != nil {
		return ctx, err
	}

	return runJob(ctx, r, j, ro)
}

func runJob(ctx context.Context, r runtime.Runtime, j *Job, ro *runOpts) (context.Context, error) {
	ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC%s] running job '%s'\n", log.ColorInfo, log.ColorNone, ro.jobName)))
	for _, step := range j.Steps {
		ctx, _, err := runStep(ctx, r, &step, ro)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}
