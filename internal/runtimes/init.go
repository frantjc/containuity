package runtimes

import (
	"context"
	"sync"

	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/docker"
)

const (
	EnvVarRuntime      = "SQNC_RUNTIME"
	DefaultRuntimeName = docker.RuntimeName
)

type NewRuntimeFunc func(context.Context) (runtime.Runtime, error)

var (
	registeredRuntimes = struct {
		sync.RWMutex
		r map[string]NewRuntimeFunc
	}{
		r: map[string]NewRuntimeFunc{},
	}
	initializedRuntimes = struct {
		sync.RWMutex
		r map[string]runtime.Runtime
	}{
		r: map[string]runtime.Runtime{},
	}
)

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

func GetDefaultRuntime(ctx context.Context) (runtime.Runtime, error) {
	return GetRuntime(ctx, DefaultRuntimeName)
}

func GetAnyRuntime(ctx context.Context) (runtime.Runtime, error) {
	return GetRuntime(ctx)
}
