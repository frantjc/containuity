package sqnc

import (
	"context"

	imageapi "github.com/frantjc/sequence/pb/v1/image"

	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) PullImage(ctx context.Context, ref string) (runtime.Image, error) {
	res, err := r.imageClient.PullImage(ctx, &imageapi.PullImageRequest{
		Ref: ref,
	})
	if err != nil {
		return nil, err
	}

	return &sqncImage{
		ref: res.Image.Ref,
	}, nil
}
