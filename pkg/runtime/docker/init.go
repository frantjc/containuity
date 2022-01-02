package docker

import (
	"context"

	dclient "github.com/docker/docker/client"
	"github.com/frantjc/sequence/pkg/runtime"
)

func init() {
	runtime.Init("docker", func(c context.Context) (runtime.Runtime, error) {
		client, err := dclient.NewClientWithOpts(dclient.FromEnv, dclient.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}

		return &dockerRuntime{client}, nil
	})
}
