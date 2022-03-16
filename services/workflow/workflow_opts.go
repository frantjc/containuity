package workflow

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

type WorkflowOpt func(s *workflowServer) error

func WithRuntime(r runtime.Runtime) WorkflowOpt {
	return func(s *workflowServer) error {
		if s.client == nil {
			s.client = &workflowClient{}
		}
		s.client.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) WorkflowOpt {
	return func(s *workflowServer) (err error) {
		if s.client == nil {
			s.client = &workflowClient{}
		}
		s.client.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime = WithRuntimeName()
