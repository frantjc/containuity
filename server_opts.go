package sequence

import (
	"context"

	"github.com/frantjc/sequence/runtime"
	_ "github.com/frantjc/sequence/runtime/docker"
)

type serverOpts struct {
	runtime runtime.Runtime
}

type ServerOpt func(*serverOpts) error

func WithRuntime(r runtime.Runtime) ServerOpt {
	return func(so *serverOpts) error {
		so.runtime = r
		return nil
	}
}

func WithRuntimeName(names ...string) ServerOpt {
	return func(so *serverOpts) (err error) {
		so.runtime, err = runtime.Get(context.Background(), names...)
		return
	}
}

var WithAnyRuntime ServerOpt = WithRuntimeName()
