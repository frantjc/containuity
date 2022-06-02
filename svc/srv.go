package svc

import (
	"net"

	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
)

type Server interface {
	Serve(net.Listener) error
}

func NewRuntimeServer(runtime runtime.Runtime) (Server, error) {
	grpcServer := grpc.NewServer()

	imageService, err := NewImageService(runtime)
	if err != nil {
		return nil, err
	}
	imageService.Register(grpcServer)

	containerService, err := NewContainerService(runtime)
	if err != nil {
		return nil, err
	}
	containerService.Register(grpcServer)

	volumeService, err := NewVolumeService(runtime)
	if err != nil {
		return nil, err
	}
	volumeService.Register(grpcServer)

	stepService, err := NewStepService(runtime)
	if err != nil {
		return nil, err
	}
	stepService.Register(grpcServer)

	jobService, err := NewJobService(runtime)
	if err != nil {
		return nil, err
	}
	jobService.Register(grpcServer)

	workflowService, err := NewWorkflowService(runtime)
	if err != nil {
		return nil, err
	}
	workflowService.Register(grpcServer)

	return grpcServer, nil
}
