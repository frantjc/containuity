package runtimes

import (
	"context"
	"os"
	"sync"

	"github.com/frantjc/sequence/runtime"
)

const (
	EnvVarRuntimeName = "SQNC_RUNTIME"
)

type NewRuntimeFunc func(context.Context) (runtime.Runtime, error)

var (
	registeredRuntimes = struct {
		sync.Mutex
		r map[string]NewRuntimeFunc
	}{
		r: map[string]NewRuntimeFunc{},
	}
	initializedRuntimes = struct {
		sync.Mutex
		r map[string]runtime.Runtime
	}{
		r: map[string]runtime.Runtime{},
	}
)

func init() {
	// TODO RegisterRuntime(sqnc.RuntimeName, ...)
}

func RegisterRuntime(name string, f NewRuntimeFunc) {
	registeredRuntimes.Lock()
	defer registeredRuntimes.Unlock()
	registeredRuntimes.r[name] = f
}

func GetRuntime(ctx context.Context, names ...string) (runtime.Runtime, error) {
	registeredRuntimes.Lock()
	defer registeredRuntimes.Unlock()

	initializedRuntimes.Lock()
	defer initializedRuntimes.Unlock()

	if len(names) == 0 {
		if name, ok := os.LookupEnv(EnvVarRuntimeName); ok {
			return GetRuntime(ctx, name)
		}

		for _, r := range initializedRuntimes.r {
			return r, nil
		}

		for _, f := range registeredRuntimes.r {
			return f(ctx)
		}
	}

	for _, name := range names {
		if r, ok := initializedRuntimes.r[name]; ok {
			return r, nil
		}

		if f, ok := registeredRuntimes.r[name]; ok {
			r, err := f(ctx)
			if err != nil {
				return nil, err
			}

			initializedRuntimes.r[name] = r

			return r, nil
		}
	}

	return nil, ErrRuntimeNotFound
}
