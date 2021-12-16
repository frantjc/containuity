package orchestrator

import (
	"context"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/runtime"
)

func RunStep(ctx context.Context, r runtime.Runtime, s *sequence.Step) error {
	return nil
}

func RunJob(ctx context.Context, r runtime.Runtime, j *sequence.Job) error {
	return nil
}

func RunWorkflow(ctx context.Context, r runtime.Runtime, j *sequence.Workflow) error {
	return nil
}
