package services

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/image"
	"github.com/frantjc/sequence/internal/convert"
	"google.golang.org/grpc"
)

func NewImageService(opts ...Opt) (ImageService, error) {
	svc := &imageServer{
		svc: &service{},
	}
	for _, opt := range opts {
		if err := opt(svc.svc); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

type imageServer struct {
	api.UnimplementedImageServer
	svc *service
}

type ImageService interface {
	api.ImageServer
	Service
}

var _ ImageService = &imageServer{}

func (s *imageServer) PullImage(ctx context.Context, in *api.PullImageRequest) (*api.PullImageResponse, error) {
	image, err := s.svc.runtime.PullImage(ctx, in.Ref)
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
