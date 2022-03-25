package workflow

import (
	"context"
	"fmt"

	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
)

func RunWorkflow(ctx context.Context, r runtime.Runtime, w *Workflow, opts ...RunOpt) error {
	ro, err := newRunOpts(append(opts, WithWorkflow(w))...)
	if err != nil {
		return err
	}

	ro.stdout.Write([]byte(fmt.Sprintf("[%sSQNC%s] running workflow '%s'\n", log.ColorInfo, log.ColorNone, w.Name)))
	for name, job := range w.Jobs {
		jobName := name
		if job.Name != "" {
			jobName = job.Name
		}

		for _, opt := range []RunOpt{
			WithJob(&job),
			WithJobName(jobName),
		} {
			err = opt(ro)
			if err != nil {
				return err
			}
		}

		if ctx, err = runJob(ctx, r, &job, ro); err != nil {
			return err
		}
	}

	return nil
}
