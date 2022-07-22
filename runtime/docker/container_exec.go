package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/sequence/runtime"
)

func (c *dockerContainer) Exec(ctx context.Context, exec *runtime.Exec, streams *runtime.Streams) error {
	execResp, err := c.client.ContainerExecCreate(ctx, c.id, types.ExecConfig{
		Cmd:          exec.GetCmd(),
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
	defer attachResp.CloseWrite() // nolint:errcheck
	defer attachResp.Close()
	_, err = stdcopy.StdCopy(streams.Out, streams.Err, attachResp.Reader)
	return err
}
