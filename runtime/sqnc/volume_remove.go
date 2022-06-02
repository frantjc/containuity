package sqnc

import (
	"context"

	volumeapi "github.com/frantjc/sequence/pb/v1/volume"
)

func (v *sqncVolume) Remove(ctx context.Context) error {
	_, err := v.client.RemoveVolume(ctx, &volumeapi.RemoveVolumeRequest{
		Name: v.source,
	})
	return err
}
