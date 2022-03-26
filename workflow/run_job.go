package workflow

import (
	"context"
	"fmt"

	"github.com/frantjc/sequence/github/actions"
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
	// ctx seems to get reset on each iteration of the loop, so even though it
	// has a ghctx at the end of one iteration, it doesn't have a ghctx at the
	// beginning of the next iteration
	//
	// as a workaround, we extract the ghctx to here
	var ghctx *actions.GlobalContext
	ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC%s] running job '%s'\n", log.ColorInfo, log.ColorNone, ro.jobName)))
	for _, step := range j.Steps {
		if ghctx != nil {
			ctx = actions.WithContext(ctx, ghctx)
		}

		ctx, _, err := runStep(ctx, r, &step, ro)
		if err != nil {
			return ctx, err
		}

		if ghctx, err = actions.Context(ctx); err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}
