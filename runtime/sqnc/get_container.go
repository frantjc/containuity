package sqnc

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/runtime"
)

func (r *sqncRuntime) GetContainer(ctx context.Context, id string) (runtime.Container, error) {
	res, err := r.runtimeClient.GetContainer(ctx, connect.NewRequest(&GetContainerRequest{
		Id: id,
	}))
	if err != nil {
		return nil, err
	}

	return &sqncContainer{
		id:     res.Msg.GetContainer().GetId(),
		client: r.runtimeClient,
	}, nil
}
