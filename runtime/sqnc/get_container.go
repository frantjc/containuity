package sqnc

import (
	"context"

	containerapi "github.com/frantjc/sequence/api/v1/container"
	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) GetContainer(ctx context.Context, id string) (runtime.Container, error) {
	res, err := r.containerClient.GetContainer(ctx, &containerapi.GetContainerRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &sqncContainer{
		id:     res.Container.Id,
		client: r.containerClient,
	}, nil
}
