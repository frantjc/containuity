package workflow

import (
	"context"

	"github.com/frantjc/sequence/log"
)

func NewWorkflowExecutor(w *Workflow, opts ...ExecOpt) (Executor, error) {
	return &workflowExecutor{w, append(opts, WithWorkflow(w))}, nil
}

type workflowExecutor struct {
	workflow *Workflow
	opts     []ExecOpt
}

var _ Executor = &workflowExecutor{}

func (e *workflowExecutor) Start(ctx context.Context) error {
	// TODO ordering, job outputs, needs, etc
	for jobName, job := range e.workflow.Jobs {
		ex, err := NewJobExecutor(job, append(e.opts, WithJobName(jobName))...)
		if err != nil {
			return err
		}

		if jex, ok := ex.(*jobExecutor); ok {
			logout := log.New(jex.stdout).SetVerbose(jex.verbose)
			logout.Infof("[%sSQNC%s] running workflow '%s'", log.ColorInfo, log.ColorNone, jex.globalContext.GitHubContext.Workflow)
		}

		if err = ex.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}
