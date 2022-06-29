package runtimes

import (
	"context"
	"os"
	"sync"

	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/docker"
)

const (
	EnvVarRuntimeName  = "SQNC_RUNTIME"
	DefaultRuntimeName = docker.RuntimeName
)

var (
	RuntimeName = DefaultRuntimeName
)

func init() {
	if runtimeName, ok := os.LookupEnv(EnvVarRuntimeName); ok {
		RuntimeName = runtimeName
	}
}

type NewRuntimeFunc func(context.Context) (runtime.Runtime, error)

var (
	RegisteredRuntimes = struct {
		sync.RWMutex
		r map[string]NewRuntimeFunc
	}{
		r: map[string]NewRuntimeFunc{},
	}
	InitializedRuntimes = struct {
		sync.RWMutex
		r map[string]runtime.Runtime
	}{
		r: map[string]runtime.Runtime{},
	}
)

func RegisterRuntime(name string, f NewRuntimeFunc) {
	RegisteredRuntimes.Lock()
	defer RegisteredRuntimes.Unlock()
	RegisteredRuntimes.r[name] = f
}

func GetRuntime(ctx context.Context, names ...string) (runtime.Runtime, error) {
	RegisteredRuntimes.Lock()
	defer RegisteredRuntimes.Unlock()

	InitializedRuntimes.Lock()
	defer InitializedRuntimes.Unlock()

	log.Infof("%s, %s", names, RegisteredRuntimes.r)

	if len(names) == 0 {
		for _, r := range InitializedRuntimes.r {
			return r, nil
		}

		for _, f := range RegisteredRuntimes.r {
			return f(ctx)
		}
	}

	for _, name := range names {
		if r, ok := InitializedRuntimes.r[name]; ok {
			return r, nil
		}

		if f, ok := RegisteredRuntimes.r[name]; ok {
			r, err := f(ctx)
			if err != nil {
				return nil, err
			}
			InitializedRuntimes.r[name] = r

			return r, nil
		}
	}

	return nil, ErrRuntimeNotFound
}

func GetDefaultRuntime(ctx context.Context) (runtime.Runtime, error) {
	return GetRuntime(ctx, RuntimeName)
}

func GetAnyRuntime(ctx context.Context) (runtime.Runtime, error) {
	return GetRuntime(ctx)
}
