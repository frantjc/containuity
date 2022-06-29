package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) CreateVolume(ctx context.Context, name string) (runtime.Volume, error) {
	res, err := r.runtimeClient.CreateVolume(ctx, connect.NewRequest(&CreateVolumeRequest{
		Name: name,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncVolume{
		source: res.Msg.GetVolume().GetSource(),
	}, nil
}
