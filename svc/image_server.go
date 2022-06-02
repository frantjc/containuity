package svc

import (
	"context"

	"github.com/frantjc/sequence/internal/convert"
	api "github.com/frantjc/sequence/pb/v1/image"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
)

func NewImageService(runtime runtime.Runtime) (ImageService, error) {
	return &imageServer{runtime: runtime}, nil
}

type imageServer struct {
	api.UnimplementedImageServer
	runtime runtime.Runtime
}

type ImageService interface {
	api.ImageServer
	Service
}

var _ ImageService = &imageServer{}

func (s *imageServer) PullImage(ctx context.Context, in *api.PullImageRequest) (*api.PullImageResponse, error) {
	image, err := s.runtime.PullImage(ctx, in.Ref)
	if err != nil {
		return nil, err
	}

	return &api.PullImageResponse{
		Image: convert.RuntimeImageToProtoImage(image),
	}, nil
}

func (s *imageServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterImageServer(r, s)
}
