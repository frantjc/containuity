package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (c *sqncContainer) Start(ctx context.Context) error {
	_, err := c.client.StartContainer(ctx, connect.NewRequest(&runtimev1.StartContainerRequest{
		Id: c.id,
	}))
	return err
}
