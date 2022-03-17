package container

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/container"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/services"
	"google.golang.org/grpc"
)

func NewService(opts ...ContainerOpt) (ContainerService, error) {
	svc := &containerServer{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

type containerServer struct {
	api.UnimplementedContainerServer
	client *containerClient
}

type ContainerService interface {
	api.ContainerServer
	services.Service
}

var _ ContainerService = &containerServer{}

func (s *containerServer) CreateContainer(ctx context.Context, in *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	return s.client.CreateContainer(ctx, in)
}

func (s *containerServer) GetContainer(ctx context.Context, in *api.GetContainerRequest) (*api.GetContainerResponse, error) {
	return s.client.GetContainer(ctx, in)
}

func (s *containerServer) ExecContainer(in *api.ExecContainerRequest, stream api.Container_ExecContainerServer) error {
	clientStream, err := s.client.ExecContainer(stream.Context(), in)
	if err != nil {
		return err
	}

	stdout, stderr := grpcio.NewLogStreamMultiplexWriter(stream)
	return grpcio.DemultiplexLogStream(clientStream, stdout, stderr)
}

func (s *containerServer) Client() (interface{}, error) {
	return s.client, nil
}

func (s *containerServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterContainerServer(r, s)
}
