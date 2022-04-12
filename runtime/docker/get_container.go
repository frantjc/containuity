package docker

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

func (r *dockerRuntime) GetContainer(ctx context.Context, id string) (runtime.Container, error) {
	return &dockerContainer{id, r.client}, nil
}
