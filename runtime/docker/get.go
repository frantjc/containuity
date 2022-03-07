package docker

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

func (r *dockerRuntime) GetContainer(ctx context.Context, id string) (runtime.Container, error) {
	return &dockerContainer{
		id:     id,
		client: r.client,
	}, nil
}
