package docker

import (
	"context"

	dclient "github.com/docker/docker/client"
	"github.com/frantjc/sequence/runtime"
)

type dockerRuntime struct {
	client *dclient.Client
}

var (
	_ runtime.Runtime = &dockerRuntime{}
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
