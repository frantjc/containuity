package runtime

import (
	"context"

	"github.com/frantjc/sequence/pkg/container"
)

type Runtime interface {
	Pull(context.Context, string) error
	Create(context.Context, *container.Spec) (container.Container, error)
}

type InitF func(context.Context) (Runtime, error)

var (
	runtimes = map[string]InitF{}
)

func Init(name string, f InitF) {
	runtimes[name] = f
}

func Get(ctx context.Context, name string) (Runtime, error) {
	if f, ok := runtimes[name]; ok {
		return f(ctx)
	}

	return nil, ErrRuntimeNotFound
}
