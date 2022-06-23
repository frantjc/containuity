package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) CreateContainer(ctx context.Context, s *runtime.Spec) (runtime.Container, error) {
	res, err := r.runtimeClient.CreateContainer(ctx, connect.NewRequest(&CreateContainerRequest{
		Spec: s,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncContainer{
		id:     res.Msg.GetContainer().GetId(),
		client: r.runtimeClient,
	}, nil
}
