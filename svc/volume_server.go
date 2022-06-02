package svc

import (
	"context"

	"github.com/frantjc/sequence/internal/convert"
	api "github.com/frantjc/sequence/pb/v1/volume"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewVolumeService(runtime runtime.Runtime) (VolumeService, error) {
	return &volumeServer{runtime: runtime}, nil
}

type volumeServer struct {
	api.UnimplementedVolumeServer
	runtime runtime.Runtime
}

type VolumeService interface {
	api.VolumeServer
	Service
}

var _ VolumeService = &volumeServer{}

func (s *volumeServer) CreateVolume(ctx context.Context, in *api.CreateVolumeRequest) (*api.CreateVolumeResponse, error) {
	volume, err := s.runtime.CreateVolume(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &api.CreateVolumeResponse{
		Volume: convert.RuntimeVolumeToProtoVolume(volume),
	}, nil
}

func (s *volumeServer) GetVolume(ctx context.Context, in *api.GetVolumeRequest) (*api.GetVolumeResponse, error) {
	volume, err := s.runtime.GetVolume(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &api.GetVolumeResponse{
		Volume: convert.RuntimeVolumeToProtoVolume(volume),
	}, nil
}

func (s *volumeServer) RemoveVolume(ctx context.Context, in *api.RemoveVolumeRequest) (*emptypb.Empty, error) {
	volume, err := s.runtime.GetVolume(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return nil, volume.Remove(ctx)
}

func (s *volumeServer) PruneVolumes(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, s.runtime.PruneVolumes(ctx)
}

func (s *volumeServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterVolumeServer(r, s)
}
