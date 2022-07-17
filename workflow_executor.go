package sequence

import (
	"context"
	"errors"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/pkg/github/actions"
)

var ErrUnmeetableJobNeeds = errors.New("job has an unmeetable needs clause")

func IsErrUnmeetableJobNeeds(err error) bool {
	return errors.Is(err, ErrUnmeetableJobNeeds)
}

type workflowExecutor struct {
	jobExecutors []*jobExecutor
	workflow     *Workflow
}

func NewWorkflowExecutor(ctx context.Context, workflow *Workflow, opts ...ExecutorOpt) (Executor, error) {
	var (
		w = &workflowExecutor{
			workflow: workflow,
		}
		needs = []string{}
		jLen  = len(workflow.Jobs)
	)
	// order Jobs so that they don't execute until after
	// their needs are met
	for len(w.jobExecutors) < jLen {
		added := false

		for jobName, job := range workflow.Jobs {
			jobOpts := opts
			jobOpts = append(jobOpts, WithWorkflowName(workflow.Name))

			if job.Needs != "" && !js.Includes(needs, job.Needs) {
				continue
			}

			if job.Name == "" {
				if js.Includes(needs, jobName) {
					continue
				}

				jobOpts = append(jobOpts, WithID(jobName), WithJobName(jobName))
			} else if js.Includes(needs, jobName) {
				continue
			}

			executor, err := NewJobExecutor(ctx, job, jobOpts...)
			if err != nil {
				return nil, err
			}

			jobExecutor := executor.(*jobExecutor)
			w.jobExecutors = append(w.jobExecutors, jobExecutor)
			needs = append(needs, jobExecutor.GlobalContext.GitHubContext.Job)
			added = true
		}

		// if we ever go a full iteration over the Jobs without
		// adding a new executor, then we have detected an infinite loop
		// due to a job having unmeetable needs
		if !added {
			return nil, ErrUnmeetableJobNeeds
		}
	}

	return w, nil
}

func (e *workflowExecutor) Execute() error {
	return e.ExecuteContext(context.Background())
}

func (e *workflowExecutor) ExecuteContext(ctx context.Context) error {
	var (
		globalContext    *actions.GlobalContext
		onWorkflowFinish Hooks[*Workflow]
	)
	for i, jobExecutor := range e.jobExecutors {
		if i != 0 {
			jobExecutor.executor.GlobalContext = globalContext
		} else {
			jobExecutor.OnWorkflowStart.Invoke(&Event[*Workflow]{
				Type:          e.workflow,
				GlobalContext: jobExecutor.GlobalContext,
			})
		}

		if err := jobExecutor.ExecuteContext(ctx); err != nil {
			return err
		}

		globalContext = jobExecutor.GlobalContext
		onWorkflowFinish = jobExecutor.OnWorkflowFinish
	}

	onWorkflowFinish.Invoke(&Event[*Workflow]{
		Type:          e.workflow,
		GlobalContext: globalContext,
	})

	return nil
}
