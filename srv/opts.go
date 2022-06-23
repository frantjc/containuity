package srv

import (
	"context"

	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/runtime"
)

type Opt func(context.Context, *Server) error

func WithRuntime(r runtime.Runtime) Opt {
	return func(ctx context.Context, s *Server) error {
		s.runtime = r
		return nil
	}
}

func WithDefaultRuntime(ctx context.Context, s *Server) error {
	var err error
	s.runtime, err = runtimes.GetDefaultRuntime(ctx)
	return err
}

func WithRuntimeName(names ...string) Opt {
	return func(ctx context.Context, s *Server) error {
		var err error
		s.runtime, err = runtimes.GetRuntime(ctx, names...)
		return err
	}
}

func WithAnyRuntime(ctx context.Context, s *Server) error {
	return WithRuntimeName()(ctx, s)
}
