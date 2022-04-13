package docker

import "context"

func (r *dockerRuntime) PruneContainers(ctx context.Context) error {
	_, err := r.client.ContainersPrune(ctx, filter)
	return err
}

func (r *dockerRuntime) PruneImages(ctx context.Context) error {
	_, err := r.client.ImagesPrune(ctx, filter)
	return err
}

func (r *dockerRuntime) PruneVolumes(ctx context.Context) error {
	_, err := r.client.VolumesPrune(ctx, filter)
	return err
}
