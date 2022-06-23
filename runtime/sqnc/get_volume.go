package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) GetVolume(ctx context.Context, name string) (runtime.Volume, error) {
	res, err := r.runtimeClient.GetVolume(ctx, connect.NewRequest(&GetVolumeRequest{
		Name: name,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncVolume{
		source: res.Msg.GetVolume().GetSource(),
		client: r.runtimeClient,
	}, nil
}
