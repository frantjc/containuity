package sqnc

import (
	"context"
	"io"

	containerapi "github.com/frantjc/sequence/api/v1/container"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
)

func (c *sqncContainer) Attach(ctx context.Context, streams *runtime.Streams) error {
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

	stream, err := c.client.AttachContainer(ctx, &containerapi.AttachContainerRequest{
		Id: c.ID(),
	})
	if err != nil {
		return err
	}

	return grpcio.DemultiplexLogStream(stream, stdout, stderr)
}
