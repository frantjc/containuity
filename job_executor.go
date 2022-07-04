package sequence

import "context"

func NewJobExecutor(ctx context.Context, job *Job, opts ...ExecutorOpt) (Executor, error) {
	internalOpts := opts
	internalOpts = append(internalOpts, func(e *executor) error {
		for k, v := range job.Env {
			e.GlobalContext.EnvContext[k] = v
		}

		return nil
	})

	if job.GetContainer() != nil && job.GetContainer().GetImage() != "" {
		internalOpts = append(internalOpts, func(e *executor) error {
			e.RunnerImage = job.Container
			return nil
		})
	}

	if job.GetName() != "" {
		internalOpts = append(internalOpts, WithID(job.GetName()))
	}

	return NewStepsExecutor(ctx, job.Steps, internalOpts...)
}
