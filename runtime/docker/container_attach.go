package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/sequence/runtime"
	"github.com/moby/term"
)

func (c *dockerContainer) Attach(ctx context.Context, streams *runtime.Streams) error {
	execResp, err := c.client.ContainerExecCreate(ctx, c.id, types.ExecConfig{
		Cmd:          []string{"/bin/sh"},
		AttachStdin:  streams.In != nil,
		AttachStdout: streams.Out != nil,
		AttachStderr: streams.Err != nil,
		DetachKeys:   streams.DetachKeys,
		Tty:          true,
	})
	if err != nil {
		return err
	}

	attachResp, err := c.client.ContainerExecAttach(ctx, execResp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer attachResp.CloseWrite() // nolint:errcheck
	defer attachResp.Close()

	detachKeysB, err := term.ToBytes(streams.DetachKeys)
	if err != nil {
		return err
	}

	errC := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(
			streams.Out,
			streams.Err,
			attachResp.Reader,
		)
		errC <- err
	}()
	go func() {
		_, err = io.Copy(attachResp.Conn, term.NewEscapeProxy(streams.In, detachKeysB))
		errC <- err
	}()

	select {
	case err := <-errC:
		if _, ok := err.(term.EscapeError); ok {
			return nil
		}

		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
