package docker

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

func (r *dockerRuntime) GetContainer(ctx context.Context, id string) (runtime.Container, error) {
	_, err := r.client.ContainerInspect(ctx, id)
	return &dockerContainer{id, r.client}, err
}
