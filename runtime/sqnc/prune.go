package sqnc

import "context"

func (r *sqncRuntime) PruneContainers(ctx context.Context) error {
	_, err := r.containerClient.PruneContainers(ctx, nil)
	return err
}

func (r *sqncRuntime) PruneImages(ctx context.Context) error {
	_, err := r.imageClient.PruneImages(ctx, nil)
	return err
}

func (r *sqncRuntime) PruneVolumes(ctx context.Context) error {
	_, err := r.volumeClient.PruneVolumes(ctx, nil)
	return err
}
