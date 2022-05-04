package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (c *dockerContainer) Remove(ctx context.Context) error {
	return c.client.ContainerRemove(ctx, c.id, types.ContainerRemoveOptions{})
}
