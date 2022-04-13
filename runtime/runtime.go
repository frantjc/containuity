package runtime

import (
	"context"
	"sync"
)

const (
	EnvVarRuntime      = "SQNC_RUNTIME"
	DockerRuntimeName  = "docker"
	DefaultRuntimeName = DockerRuntimeName
)

type Runtime interface {
	PullImage(context.Context, string) (Image, error)
	CreateContainer(context.Context, *Spec) (Container, error)
	GetContainer(context.Context, string) (Container, error)
	CreateVolume(context.Context, string) (Volume, error)
	GetVolume(context.Context, string) (Volume, error)
	PruneContainers(context.Context) error
	PruneImages(context.Context) error
	PruneVolumes(context.Context) error
}

type InitF func(context.Context) (Runtime, error)

var (
	registeredRuntimes = struct {
		sync.RWMutex
		r map[string]InitF
	}{
		r: map[string]InitF{},
	}
	initializedRuntimes = struct {
		sync.RWMutex
		r map[string]Runtime
	}{
		r: map[string]Runtime{},
	}
)

func Init(name string, f InitF) {
	registeredRuntimes.Lock()
	defer registeredRuntimes.Unlock()
	registeredRuntimes.r[name] = f
}

func Get(ctx context.Context, names ...string) (Runtime, error) {
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
