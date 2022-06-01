package svc

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"

	api "github.com/frantjc/sequence/api/v1/container"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewContainerService(runtime runtime.Runtime) (ContainerService, error) {
	return &containerServer{runtime: runtime}, nil
}

type containerServer struct {
	api.UnimplementedContainerServer
	runtime runtime.Runtime
}

type ContainerService interface {
	api.ContainerServer
	Service
}

var _ ContainerService = &containerServer{}

func (s *containerServer) CreateContainer(ctx context.Context, in *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	container, err := s.runtime.CreateContainer(ctx, convert.ProtoSpecToRuntimeSpec(in.Spec))
	if err != nil {
		return nil, err
	}

	return &api.CreateContainerResponse{
		Container: convert.RuntimeContainerToProtoContainer(container),
	}, nil
}

func (s *containerServer) GetContainer(ctx context.Context, in *api.GetContainerRequest) (*api.GetContainerResponse, error) {
	container, err := s.runtime.GetContainer(ctx, in.Id)
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
	container, err := s.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return err
	}

	return container.Exec(ctx, runtime.NewStreams(
		os.Stdin,
		stdout,
		stderr,
	))
}

func (s *containerServer) StartContainer(ctx context.Context, in *api.StartContainerRequest) (*emptypb.Empty, error) {
	container, err := s.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return nil, container.Start(ctx)
}

func (s *containerServer) AttachContainer(in *api.AttachContainerRequest, stream api.Container_AttachContainerServer) error {
	var (
		ctx            = stream.Context()
		stdout, stderr = grpcio.NewLogStreamMultiplexWriter(stream)
	)
	container, err := s.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return err
	}

	return container.Attach(ctx, runtime.NewStreams(
		os.Stdin,
		stdout,
		stderr,
	))
}

func (s *containerServer) RemoveContainer(ctx context.Context, in *api.RemoveContainerRequest) (*emptypb.Empty, error) {
	container, err := s.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return nil, container.Remove(ctx)
}

func (s *containerServer) PruneContainers(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, s.runtime.PruneContainers(ctx)
}

func (s *containerServer) CopyToContainer(ctx context.Context, in *api.CopyToContainerRequest) (*emptypb.Empty, error) {
	container, err := s.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return nil, container.CopyTo(ctx, bytes.NewReader(in.Content), in.Destination)
}

func (s *containerServer) CopyFromContainer(ctx context.Context, in *api.CopyFromContainerRequest) (*api.CopyFromContainerResponse, error) {
	container, err := s.runtime.GetContainer(ctx, in.Id)
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
