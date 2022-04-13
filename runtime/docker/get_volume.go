package docker

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

func (r *dockerRuntime) GetVolume(ctx context.Context, name string) (runtime.Volume, error) {
	_, err := r.client.VolumeInspect(ctx, name)
	return &dockerVolume{name, r.client}, err
}
