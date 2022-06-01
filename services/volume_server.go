package services

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/volume"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewVolumeService(runtime runtime.Runtime) (VolumeService, error) {
	svc := &volumeServer{
		svc: &service{runtime},
	}
	return svc, nil
}

type volumeServer struct {
	api.UnimplementedVolumeServer
	svc *service
}

type VolumeService interface {
	api.VolumeServer
	Service
}

var _ VolumeService = &volumeServer{}

func (s *volumeServer) CreateVolume(ctx context.Context, in *api.CreateVolumeRequest) (*api.CreateVolumeResponse, error) {
	volume, err := s.svc.runtime.CreateVolume(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &api.CreateVolumeResponse{
		Volume: convert.RuntimeVolumeToProtoVolume(volume),
	}, nil
}

func (s *volumeServer) GetVolume(ctx context.Context, in *api.GetVolumeRequest) (*api.GetVolumeResponse, error) {
	volume, err := s.svc.runtime.GetVolume(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &api.GetVolumeResponse{
		Volume: convert.RuntimeVolumeToProtoVolume(volume),
	}, nil
}

func (s *volumeServer) RemoveVolume(ctx context.Context, in *api.RemoveVolumeRequest) (*emptypb.Empty, error) {
	volume, err := s.svc.runtime.GetVolume(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return nil, volume.Remove(ctx)
}

func (s *volumeServer) PruneVolumes(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, s.svc.runtime.PruneVolumes(ctx)
}

func (s *volumeServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterVolumeServer(r, s)
}
