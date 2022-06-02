package sqnc

import (
	"context"

	containerapi "github.com/frantjc/sequence/pb/v1/container"
)

func (c *sqncContainer) Remove(ctx context.Context) error {
	_, err := c.client.RemoveContainer(ctx, &containerapi.RemoveContainerRequest{
		Id: c.id,
	})
	return err
}
