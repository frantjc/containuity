package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/sequence/runtime"
)

func (c *dockerContainer) Exec(ctx context.Context, streams *runtime.Streams) error {
	attachResp, err := c.client.ContainerAttach(ctx, c.id, types.ContainerAttachOptions{
		Stream: true,
		Stdout: streams.Stdout != nil,
		Stderr: streams.Stderr != nil,
	})
	if err != nil {
		return err
	}
	go stdcopy.StdCopy(streams.Stdout, streams.Stderr, attachResp.Reader) //nolint:errcheck

	err = c.client.ContainerStart(ctx, c.id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusC, errC := c.client.ContainerWait(ctx, c.id, container.WaitConditionNotRunning)
	select {
	case err := <-errC:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return ctx.Err()
	case status := <-statusC:
		if status.StatusCode != 0 {
			msg := ""
			if status.Error != nil {
				msg = fmt.Sprintf(": %s", status.Error.Message)
			}
			return fmt.Errorf("container exited with nonzero code %d%s", status.StatusCode, msg)
		}
	}

	return nil
}
