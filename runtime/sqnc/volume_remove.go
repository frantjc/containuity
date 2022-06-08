package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (v *sqncVolume) Remove(ctx context.Context) error {
	_, err := v.client.RemoveVolume(ctx, connect.NewRequest(&runtimev1.RemoveVolumeRequest{
		Name: v.source,
	}))
	return err
}
