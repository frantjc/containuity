package sqnc

import "context"

func (r *sqncRuntime) PruneContainers(ctx context.Context) error {
	_, err := r.runtimeClient.PruneContainers(ctx, nil)
	return err
}

func (r *sqncRuntime) PruneImages(ctx context.Context) error {
	_, err := r.runtimeClient.PruneImages(ctx, nil)
	return err
}

func (r *sqncRuntime) PruneVolumes(ctx context.Context) error {
	_, err := r.runtimeClient.PruneVolumes(ctx, nil)
	return err
}
