package sqnc

import (
	"context"

	containerapi "github.com/frantjc/sequence/pb/v1/container"
)

func (c *sqncContainer) Start(ctx context.Context) error {
	_, err := c.client.StartContainer(ctx, &containerapi.StartContainerRequest{
		Id: c.id,
	})
	return err
}
