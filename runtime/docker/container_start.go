package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (c *dockerContainer) Start(ctx context.Context) error {
	return c.client.ContainerStart(ctx, c.id, types.ContainerStartOptions{})
}
