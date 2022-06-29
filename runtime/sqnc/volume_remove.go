package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
)

func (v *sqncVolume) Remove(ctx context.Context) error {
	_, err := v.client.RemoveVolume(ctx, connect.NewRequest(&RemoveVolumeRequest{
		Name: v.source,
	}))
	return err
}
