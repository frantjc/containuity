package docker

import "context"

func (v *dockerVolume) Remove(ctx context.Context) error {
	return v.client.VolumeRemove(ctx, v.name, true)
}
