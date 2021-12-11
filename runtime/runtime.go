package runtime

import (
	"context"
	"fmt"

	"github.com/frantjc/sequence"
	"github.com/google/uuid"
)

type runOpts struct {
	id string
}

func defaultRunOpts() *runOpts {
	return &runOpts{
		id: uuid.NewString(),
	}
}

func createRunOpts(opts ...RunOpt) (*runOpts, error) {
	ropts := defaultRunOpts()
	for _, opt := range opts {
		err := opt(ropts)
		if err != nil {
			return nil, err
		}
	}

	return ropts, nil
}

type RunOpt func(r *runOpts) error

type Runtime interface {
	Run(context.Context, *sequence.Step, ...RunOpt) error
}

type initRuntime func() (Runtime, error)

var (
	runtimes = map[string]initRuntime{}
)

func RegisterRuntime(name string, f initRuntime) {
	runtimes[name] = f
}

func GetRuntime(name string) (Runtime, error) {
	if f, ok := runtimes[name]; ok {
		return f()
	}

	return nil, fmt.Errorf("unknown runtime: %s", name)
}
