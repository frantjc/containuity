package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
)

func (c *sqncContainer) Remove(ctx context.Context) error {
	_, err := c.client.RemoveContainer(ctx, connect.NewRequest(&RemoveContainerRequest{
		Id: c.id,
	}))
	return err
}
