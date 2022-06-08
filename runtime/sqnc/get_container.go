package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (r *sqncRuntime) GetContainer(ctx context.Context, id string) (runtime.Container, error) {
	res, err := r.containerClient.GetContainer(ctx, connect.NewRequest(&runtimev1.GetContainerRequest{
		Id: id,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncContainer{
		id:     res.Msg.GetContainer().GetId(),
		client: r.containerClient,
	}, nil
}
