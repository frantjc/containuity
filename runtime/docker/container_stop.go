package docker

import "context"

func (c *dockerContainer) Stop(ctx context.Context) error {
	return c.client.ContainerStop(ctx, c.id, nil)
}
