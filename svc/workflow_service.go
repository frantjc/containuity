package svc

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence"
)

type WorkflowServiceHandler struct {
	sequence.UnimplementedWorkflowServiceHandler
}

func (*WorkflowServiceHandler) RunWorkflow(context.Context, *connect.Request[sequence.RunWorkflowRequest], *connect.ServerStream[sequence.RunWorkflowResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("sequence.v1.WorkflowService.RunWorkflow is not implemented"))
}

func (*WorkflowServiceHandler) RunJob(context.Context, *connect.Request[sequence.RunJobRequest], *connect.ServerStream[sequence.RunJobResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("sequence.v1.WorkflowService.RunJob is not implemented"))
}

func (*WorkflowServiceHandler) RunStep(context.Context, *connect.Request[sequence.RunStepRequest], *connect.ServerStream[sequence.RunStepResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("sequence.v1.WorkflowService.RunStep is not implemented"))
}
