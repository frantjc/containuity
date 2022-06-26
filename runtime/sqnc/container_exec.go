package sqnc

import (
	"context"
	"fmt"
	"io"

	"github.com/bufbuild/connect-go"
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

	_, err := c.client.ExecContainer(ctx, connect.NewRequest(&ExecContainerRequest{
		Id: c.GetID(),
	}))
	if err != nil {
		return err
	}

	var _ = stderr

	return fmt.Errorf("unimplemented")
}