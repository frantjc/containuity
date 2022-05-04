package services

import (
	"context"

	"github.com/frantjc/sequence/api/types"
	api "github.com/frantjc/sequence/api/v1/volume"
	"github.com/frantjc/sequence/internal/convert"
	"google.golang.org/grpc"
)

func NewVolumeService(opts ...Opt) (VolumeService, error) {
	svc := &volumeServer{
		svc: &service{},
	}
	for _, opt := range opts {
		if err := opt(svc.svc); err != nil {
			return nil, err
		}
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

func (s *volumeServer) RemoveVolume(ctx context.Context, in *api.RemoveVolumeRequest) (*types.Empty, error) {
	volume, err := s.svc.runtime.GetVolume(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &types.Empty{}, volume.Remove(ctx)
}

func (s *volumeServer) PruneVolumes(ctx context.Context, _ *types.Empty) (*types.Empty, error) {
	return &types.Empty{}, s.svc.runtime.PruneVolumes(ctx)
}

func (s *volumeServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterVolumeServer(r, s)
}
