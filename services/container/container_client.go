package container

import (
	"context"
	"os"

	api "github.com/frantjc/sequence/api/v1/container"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
)

type containerClient struct {
	runtime runtime.Runtime
}

var _ api.ContainerClient = &containerClient{}

func (c *containerClient) CreateContainer(ctx context.Context, req *api.CreateContainerRequest, _ ...grpc.CallOption) (*api.CreateContainerResponse, error) {
	container, err := c.runtime.CreateContainer(ctx, convert.ProtoSpecToRuntimeSpec(req.Spec))
	if err != nil {
		return nil, err
	}

	return &api.CreateContainerResponse{
		Container: convert.RuntimeContainerToProtoContainer(container),
	}, nil
}

func (c *containerClient) GetContainer(ctx context.Context, req *api.GetContainerRequest, _ ...grpc.CallOption) (*api.GetContainerResponse, error) {
	container, err := c.runtime.GetContainer(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &api.GetContainerResponse{
		Container: convert.RuntimeContainerToProtoContainer(container),
	}, nil
}

func (c *containerClient) ExecContainer(ctx context.Context, in *api.ExecContainerRequest, _ ...grpc.CallOption) (api.Container_ExecContainerClient, error) {
	var (
		stream         = grpcio.NewLogStream(ctx)
		stdout, stderr = grpcio.NewLogStreamMultiplexWriter(stream)
	)
	container, err := c.runtime.GetContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	go func() {
		defer stream.CloseSend()
		if err := container.Exec(ctx, runtime.ExecStreams(
			os.Stdin,
			stdout,
			stderr,
		)); err != nil {
			stream.SendErr(err)
		}
	}()

	return stream, nil
}
