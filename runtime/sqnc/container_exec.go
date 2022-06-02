package sqnc

import (
	"context"
	"io"

	"github.com/frantjc/sequence/internal/grpcio"
	containerapi "github.com/frantjc/sequence/pb/v1/container"
	"github.com/frantjc/sequence/runtime"
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

	stream, err := c.client.ExecContainer(ctx, &containerapi.ExecContainerRequest{
		Id: c.ID(),
	})
	if err != nil {
		return err
	}

	return grpcio.DemultiplexLogStream(stream, stdout, stderr)
}
