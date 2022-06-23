package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) PullImage(ctx context.Context, ref string) (runtime.Image, error) {
	res, err := r.runtimeClient.PullImage(ctx, connect.NewRequest(&PullImageRequest{
		Ref: ref,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncImage{
		ref: res.Msg.GetImage().GetRef(),
	}, nil
}
