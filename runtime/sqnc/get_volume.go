package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (r *sqncRuntime) GetVolume(ctx context.Context, name string) (runtime.Volume, error) {
	res, err := r.volumeClient.GetVolume(ctx, connect.NewRequest(&runtimev1.GetVolumeRequest{
		Name: name,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncVolume{
		source: res.Msg.GetVolume().GetSource(),
		client: r.volumeClient,
	}, nil
}
