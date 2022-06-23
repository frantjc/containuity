package docker

import (
	"context"

	dclient "github.com/docker/docker/client"
	"github.com/frantjc/sequence/runtime"
)

func NewRuntime(ctx context.Context) (runtime.Runtime, error) {
	client, err := dclient.NewClientWithOpts(dclient.FromEnv, dclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &dockerRuntime{client}, nil
}

const RuntimeName = "docker"

type dockerRuntime struct {
	client *dclient.Client
}

var _ runtime.Runtime = &dockerRuntime{}
