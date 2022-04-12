package docker

import (
	"context"

	"github.com/docker/docker/api/types/volume"
	"github.com/frantjc/sequence/runtime"
)

func (r *dockerRuntime) CreateVolume(ctx context.Context, name string) (runtime.Volume, error) {
	volume, err := r.client.VolumeCreate(ctx, volume.VolumeCreateBody{
		Driver: "local",
		Labels: labels,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	return &dockerVolume{volume.Name, r.client}, nil
}
