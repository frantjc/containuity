package srv

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

type Opt func(context.Context, *Server) error

func WithRuntime(r runtime.Runtime) Opt {
	return func(ctx context.Context, s *Server) error {
		s.runtime = r
		return nil
	}
}
