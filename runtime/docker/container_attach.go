package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/frantjc/sequence/runtime"
)

func (c *dockerContainer) Attach(ctx context.Context, streams *runtime.Streams) error {
	execResp, err := c.client.ContainerExecCreate(ctx, c.id, types.ExecConfig{
		Cmd:          []string{"sh"},
		AttachStdin:  streams.In != nil,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}

	attachResp, err := c.client.ContainerExecAttach(ctx, execResp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer attachResp.Close()

	streamer := &hijackedIOStreamer{
		inputStream:  io.NopCloser(streams.In),
		outputStream: streams.Out,
		errorStream:  streams.Err,
		resp:         attachResp,
	}

	return streamer.stream(ctx)
}
