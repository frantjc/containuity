package sqnc

import (
	"context"
	"io"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/internal/protobufio"
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (c *sqncContainer) Exec(ctx context.Context, streams *runtime.Streams) error {
	var (
		stdout = io.Discard
		stderr = stdout
	)
	if streams.Stdout != nil {
		stdout = streams.Stdout
	}
	if streams.Stderr != nil {
		stderr = streams.Stderr
	}

	stream, err := c.client.ExecContainer(ctx, connect.NewRequest(&runtimev1.ExecContainerRequest{
		Id: c.GetID(),
	}))
	if err != nil {
		return err
	}

	return protobufio.DemultiplexLogStream[*runtimev1.ExecContainerResponse](stream, stdout, stderr)
}
