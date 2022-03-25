package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/sequence/runtime"
)

type dockerContainer struct {
	id     string
	client *client.Client
}

var (
	_ runtime.Container = &dockerContainer{}
)

func (c *dockerContainer) ID() string {
	return c.id
}

func (c *dockerContainer) Exec(ctx context.Context, opts ...runtime.ExecOpt) error {
	e, err := runtime.NewExec(opts...)
	if err != nil {
		return err
	}

	attachResp, err := c.client.ContainerAttach(ctx, c.id, types.ContainerAttachOptions{
		Stream: true,
		Stdout: e.Stdout != nil,
		Stderr: e.Stderr != nil,
	})
	if err != nil {
		return err
	}
	go stdcopy.StdCopy(e.Stdout, e.Stderr, attachResp.Reader)

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
		defer c.client.ContainerRemove(ctx, c.id, types.ContainerRemoveOptions{
			Force: true,
		})
		return ctx.Err()
	case <-statusC:
	}

	return nil
}
