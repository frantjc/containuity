package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (r *sqncRuntime) CreateContainer(ctx context.Context, s *runtimev1.Spec) (runtime.Container, error) {
	res, err := r.containerClient.CreateContainer(ctx, connect.NewRequest(&runtimev1.CreateContainerRequest{
		Spec: s,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncContainer{
		id:     res.Msg.GetContainer().GetId(),
		client: r.containerClient,
	}, nil
}
