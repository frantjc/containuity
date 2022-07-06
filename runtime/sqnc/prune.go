package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
)

func (r *sqncRuntime) PruneContainers(ctx context.Context) error {
	_, err := r.runtimeClient.PruneContainers(ctx, connect.NewRequest(&PruneContainersRequest{}))
	return err
}

func (r *sqncRuntime) PruneImages(ctx context.Context) error {
	_, err := r.runtimeClient.PruneImages(ctx, connect.NewRequest(&PruneImagesRequest{}))
	return err
}

func (r *sqncRuntime) PruneVolumes(ctx context.Context) error {
	_, err := r.runtimeClient.PruneVolumes(ctx, connect.NewRequest(&PruneVolumesRequest{}))
	return err
}
