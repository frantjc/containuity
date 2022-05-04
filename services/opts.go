package services

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

type service struct {
	runtime runtime.Runtime
}

type Opt func(*service) error

func WithRuntime(r runtime.Runtime) Opt {
	return func(s *service) error {
		s.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) Opt {
	return func(s *service) (err error) {
		s.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime = WithRuntimeName()
