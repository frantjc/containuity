package sqnc

import (
	"context"

	containerapi "github.com/frantjc/sequence/api/v1/container"

	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) CreateContainer(ctx context.Context, s *runtime.Spec) (runtime.Container, error) {
	res, err := r.containerClient.CreateContainer(ctx, &containerapi.CreateContainerRequest{
		Spec: convert.RuntimeSpecToProtoSpec(s),
	})
	if err != nil {
		return nil, err
	}

	return &sqncContainer{
		id:     res.Container.Id,
		client: r.containerClient,
	}, nil
}
