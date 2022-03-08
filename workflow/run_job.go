package workflow

import (
	"context"
	"fmt"

	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
)

func RunJob(ctx context.Context, r runtime.Runtime, j *Job, opts ...RunOpt) error {
	ro, err := newRunOpts(append(opts, WithJob(j))...)
	if err != nil {
		return err
	}

	return runJob(ctx, r, j, ro)
}

func runJob(ctx context.Context, r runtime.Runtime, j *Job, ro *runOpts) error {
	ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC%s] running job %s\n", log.ColorInfo, log.ColorNone, ro.jobName)))
	for _, step := range j.Steps {
		err := runStep(ctx, r, &step, ro)
		if err != nil {
			return err
		}
	}

	return nil
}
