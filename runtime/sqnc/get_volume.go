package sqnc

import (
	"context"

	volumeapi "github.com/frantjc/sequence/api/v1/volume"

	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) GetVolume(ctx context.Context, name string) (runtime.Volume, error) {
	res, err := r.volumeClient.GetVolume(ctx, &volumeapi.GetVolumeRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return &sqncVolume{
		source: res.Volume.Source,
		client: r.volumeClient,
	}, nil
}
