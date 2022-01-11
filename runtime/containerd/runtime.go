package containerd

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/frantjc/sequence/runtime"
)

type containerdRuntime struct {
	client *containerd.Client
}

var (
	_ runtime.Runtime = &containerdRuntime{}
)

func init() {
	runtime.Init("containerd", func(c context.Context) (runtime.Runtime, error) {
		client, err := containerd.New("/run/containerd/containerd.sock")
		if err != nil {
			return nil, err
		}

		return &containerdRuntime{client}, nil
	})
}
