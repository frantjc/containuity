package runtime

import (
	"context"
	"fmt"

	dtypes "github.com/docker/docker/api/types"
	dcontainer "github.com/docker/docker/api/types/container"
	dmount "github.com/docker/docker/api/types/mount"
	dstdcopy "github.com/docker/docker/pkg/stdcopy"

	"github.com/frantjc/sequence/internal/image"
	"github.com/frantjc/sequence/pkg/container"
	"github.com/frantjc/sequence/pkg/io"
)

type dockerContainer struct {
	id string
	r  *dockerRuntime
}

func (r *dockerRuntime) Create(ctx context.Context, c *container.Spec) (container.Container, error) {
	if r == nil {
		return nil, fmt.Errorf("nil runtime")
	}

	ref, err := image.ParseRef(c.Image)
	if err != nil {
		return nil, err
	}

	conf := &dcontainer.Config{
		Image:      ref,
		Entrypoint: c.Entrypoint,
		Cmd:        c.Cmd,
		Env:        c.Env,
		Labels:     defaultLabels(),
	}

	hconf := &dcontainer.HostConfig{
		AutoRemove: true,
		Privileged: c.Privileged,
		Mounts:     defaultMounts(),
	}

	for _, m := range c.Mounts {
		dt := dmount.Type(m.Type)
		dm := dmount.Mount{
			Type:   dt,
			Source: m.Source,
			Target: m.Destination,
		}

		switch dt {
		case dmount.TypeBind:
		case dmount.TypeVolume:
			dm.VolumeOptions.Labels = defaultLabels()
		case dmount.TypeTmpfs:
			dm.Source = ""
		}

		for _, opt := range m.Options {
			switch opt {
			case container.MountOptReadOnly:
				dm.ReadOnly = true
			}
		}

		hconf.Mounts = append(hconf.Mounts, dm)
	}

	createResp, err := r.client.ContainerCreate(ctx, conf, hconf, nil, nil, "")
	if err != nil {
		return nil, err
	}
	id := createResp.ID

	return &dockerContainer{id, r}, nil
}

func (c *dockerContainer) Start(ctx context.Context, s io.Streams) error {
	attachResp, err := c.r.client.ContainerAttach(ctx, c.id, dtypes.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	// TODO should we do this?
	// go io.Copy(attachResp.Conn, c.opts.Stdin)
	// TODO prettify this
	go dstdcopy.StdCopy(s.Stdout, s.Stderr, attachResp.Reader)

	err = c.r.client.ContainerStart(ctx, c.id, dtypes.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusC, errC := c.r.client.ContainerWait(ctx, c.id, dcontainer.WaitConditionNotRunning)
	select {
	case err := <-errC:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		defer c.r.client.ContainerRemove(ctx, c.id, dtypes.ContainerRemoveOptions{
			Force: true,
		})
		return ctx.Err()
	case <-statusC:
	}

	return nil
}
