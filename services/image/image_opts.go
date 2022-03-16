package image

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

type ImageOpt func(s *imageServer) error

func WithRuntime(r runtime.Runtime) ImageOpt {
	return func(s *imageServer) error {
		if s.client == nil {
			s.client = &imageClient{}
		}
		s.client.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) ImageOpt {
	return func(s *imageServer) (err error) {
		if s.client == nil {
			s.client = &imageClient{}
		}
		s.client.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime = WithRuntimeName()
