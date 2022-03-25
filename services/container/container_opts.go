package container

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

type ContainerOpt func(*containerServer) error

func WithRuntime(r runtime.Runtime) ContainerOpt {
	return func(s *containerServer) error {
		if s.client == nil {
			s.client = &containerClient{}
		}
		s.client.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) ContainerOpt {
	return func(s *containerServer) (err error) {
		if s.client == nil {
			s.client = &containerClient{}
		}
		s.client.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime = WithRuntimeName()
