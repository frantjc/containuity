package step

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

type StepOpt func(s *stepServer) error

func WithRuntime(r runtime.Runtime) StepOpt {
	return func(s *stepServer) error {
		if s.client == nil {
			s.client = &stepClient{}
		}
		s.client.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) StepOpt {
	return func(s *stepServer) (err error) {
		if s.client == nil {
			s.client = &stepClient{}
		}
		s.client.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime = WithRuntimeName()
