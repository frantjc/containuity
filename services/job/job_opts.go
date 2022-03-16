package job

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

type JobOpt func(s *jobServer) error

func WithRuntime(r runtime.Runtime) JobOpt {
	return func(s *jobServer) error {
		if s.client == nil {
			s.client = &jobClient{}
		}
		s.client.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) JobOpt {
	return func(s *jobServer) (err error) {
		if s.client == nil {
			s.client = &jobClient{}
		}
		s.client.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime = WithRuntimeName()
