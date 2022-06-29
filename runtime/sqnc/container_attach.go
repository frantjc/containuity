package sqnc

import (
	"context"
	"fmt"
	"io"

	"github.com/bufbuild/connect-go"
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

	_, err := c.client.AttachContainer(ctx, connect.NewRequest(&AttachContainerRequest{
		Id: c.GetID(),
	}))
	if err != nil {
		return err
	}

	var (
		_ = stdout
		_ = stderr
	)

	// TODO
	return fmt.Errorf("unimplemented")
}
