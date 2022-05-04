package services

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"

	"github.com/frantjc/sequence/api/types"
	api "github.com/frantjc/sequence/api/v1/container"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
)

func NewContainerService(opts ...Opt) (ContainerService, error) {
	svc := &containerServer{
		svc: &service{},
	}
	for _, opt := range opts {
		if err := opt(svc.svc); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

type containerServer struct {
	api.UnimplementedContainerServer
	svc *service
}

type ContainerService interface {
	api.ContainerServer
	Service
}

var _ ContainerService = &containerServer{}

func (s *containerServer) CreateContainer(ctx context.Context, in *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	container, err := s.svc.runtime.CreateContainer(ctx, convert.ProtoSpecToRuntimeSpec(in.Spec))
	if err != nil {
		return nil, err
	}

	return &api.CreateContainerResponse{
		Container: convert.RuntimeContainerToProtoContainer(container),
	}, nil
}

func (s *containerServer) GetContainer(ctx context.Context, in *api.GetContainerRequest) (*api.GetContainerResponse, error) {
	container, err := s.svc.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &api.GetContainerResponse{
		Container: convert.RuntimeContainerToProtoContainer(container),
	}, nil
}

func (s *containerServer) ExecContainer(in *api.ExecContainerRequest, stream api.Container_ExecContainerServer) error {
	var (
		ctx            = stream.Context()
		stdout, stderr = grpcio.NewLogStreamMultiplexWriter(stream)
	)
	container, err := s.svc.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return err
	}

	return container.Exec(ctx, runtime.NewStreams(
		os.Stdin,
		stdout,
		stderr,
	))
}

func (s *containerServer) StartContainer(ctx context.Context, in *api.StartContainerRequest) (*types.Empty, error) {
	container, err := s.svc.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &types.Empty{}, container.Start(ctx)
}

func (s *containerServer) AttachContainer(in *api.AttachContainerRequest, stream api.Container_AttachContainerServer) error {
	var (
		ctx            = stream.Context()
		stdout, stderr = grpcio.NewLogStreamMultiplexWriter(stream)
	)
	container, err := s.svc.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return err
	}

	return container.Attach(ctx, runtime.NewStreams(
		os.Stdin,
		stdout,
		stderr,
	))
}

func (s *containerServer) RemoveContainer(ctx context.Context, in *api.RemoveContainerRequest) (*types.Empty, error) {
	container, err := s.svc.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &types.Empty{}, container.Remove(ctx)
}

func (s *containerServer) PruneContainers(ctx context.Context, in *types.Empty) (*types.Empty, error) {
	return &types.Empty{}, s.svc.runtime.PruneContainers(ctx)
}

func (s *containerServer) CopyToContainer(ctx context.Context, in *api.CopyToContainerRequest) (*types.Empty, error) {
	container, err := s.svc.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &types.Empty{}, container.CopyTo(ctx, bytes.NewReader(in.Content), in.Destination)
}

func (s *containerServer) CopyFromContainer(ctx context.Context, in *api.CopyFromContainerRequest) (*api.CopyFromContainerResponse, error) {
	container, err := s.svc.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	content, err := container.CopyFrom(ctx, in.Source)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, err
	}

	return &api.CopyFromContainerResponse{
		Content: b,
	}, nil
}

func (s *containerServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterContainerServer(r, s)
}
