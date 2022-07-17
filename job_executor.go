package sequence

import (
	"context"

	"github.com/frantjc/sequence/pkg/github/actions"
)

type jobExecutor struct {
	*stepsExecutor
	job *Job
}

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
		internalOpts = append(internalOpts, WithID(job.GetName()), WithJobName(job.GetName()))
	}

	executor, err := NewStepsExecutor(ctx, job.Steps, internalOpts...)
	if err != nil {
		return nil, err
	}

	return &jobExecutor{
		stepsExecutor: executor.(*stepsExecutor),
		job:           job,
	}, nil
}

func (e *jobExecutor) Execute() error {
	return e.ExecuteContext(context.Background())
}

func (e *jobExecutor) ExecuteContext(ctx context.Context) error {
	e.OnJobStart.Invoke(&Event[*Job]{
		Type:          e.job,
		GlobalContext: e.GlobalContext,
	})

	if err := e.stepsExecutor.ExecuteContext(ctx); err != nil {
		return err
	}

	if e.ID != "" {
		e.stepsExecutor.GlobalContext.NeedsContext[e.ID] = &actions.NeedsContext{
			Outputs: map[string]string{},
		}

		expander := actions.NewExpander(e.GlobalContext.GetString)
		for k, v := range e.job.Outputs {
			e.stepsExecutor.GlobalContext.NeedsContext[e.ID].Outputs[k] = expander.Expand(v)
		}
	}

	e.OnJobFinish.Invoke(&Event[*Job]{
		Type:          e.job,
		GlobalContext: e.GlobalContext,
	})

	return nil
}
