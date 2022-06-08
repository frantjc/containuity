package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (r *sqncRuntime) PullImage(ctx context.Context, ref string) (runtime.Image, error) {
	res, err := r.imageClient.PullImage(ctx, connect.NewRequest(&runtimev1.PullImageRequest{
		Ref: ref,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncImage{
		ref: res.Msg.GetImage().GetRef(),
	}, nil
}
