package sequence

import (
	"context"
	"fmt"
)

type Runtime interface {
	Run(context.Context, Steppable) error
}

type initRuntime func(context.Context) (Runtime, error)

var runtimes map[string]initRuntime

func init() {
	runtimes = map[string]initRuntime{}
}

func RegisterRuntime(name string, f initRuntime) {
	runtimes[name] = f
}

func GetRuntime(ctx context.Context, name string) (Runtime, error) {
	if f, ok := runtimes[name]; ok {
		return f(ctx)
	}

	return nil, fmt.Errorf("unknown runtime: %s", name)
}
