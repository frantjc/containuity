package image

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/image"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
)

type imageClient struct {
	runtime runtime.Runtime
}

var _ api.ImageClient = &imageClient{}

func (c *imageClient) PullImage(ctx context.Context, in *api.PullImageRequest, _ ...grpc.CallOption) (*api.PullImageResponse, error) {
	image, err := c.runtime.PullImage(ctx, in.Ref)
	if err != nil {
		return nil, err
	}

	return &api.PullImageResponse{
		Image: convert.RuntimeImageToProtoImage(image),
	}, nil
}
