package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/filters"
)

func (r *dockerRuntime) Prune(ctx context.Context) error {
	filter := filters.NewArgs()
	for k, v := range labels {
		filter.Add("label", fmt.Sprintf("%s=%s", k, v))
	}

	_, err := r.client.ContainersPrune(ctx, filter)
	if err != nil {
		return err
	}

	_, err = r.client.VolumesPrune(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
