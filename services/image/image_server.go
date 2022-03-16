package image

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/image"
	"github.com/frantjc/sequence/services"
	"google.golang.org/grpc"
)

func NewService(opts ...ImageOpt) (ImageService, error) {
	svc := &imageServer{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

type imageServer struct {
	api.UnimplementedImageServer
	client *imageClient
}

type ImageService interface {
	api.ImageServer
	services.Service
}

var _ ImageService = &imageServer{}

func (s *imageServer) PullImage(ctx context.Context, in *api.PullImageRequest) (*api.PullImageResponse, error) {
	return s.client.PullImage(ctx, in)
}

func (s *imageServer) Client() (interface{}, error) {
	return s.client, nil
}

func (s *imageServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterImageServer(r, s)
}
