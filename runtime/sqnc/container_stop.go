package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
)

func (c *sqncContainer) Stop(ctx context.Context) error {
	_, err := c.client.StopContainer(ctx, connect.NewRequest(&StopContainerRequest{
		Id: c.id,
	}))
	return err
}
