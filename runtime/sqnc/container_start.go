package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
)

func (c *sqncContainer) Start(ctx context.Context) error {
	_, err := c.client.StartContainer(ctx, connect.NewRequest(&StartContainerRequest{
		Id: c.id,
	}))
	return err
}
