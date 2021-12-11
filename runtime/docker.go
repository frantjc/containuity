package runtime

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/sequence"
)

var labels = map[string]string{
	"sequence": "true",
}

func init() {
	RegisterRuntime("docker", func() (Runtime, error) {
		c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}

		return &dockerRuntime{client: c}, nil
	})
}

type dockerRuntime struct {
	client *client.Client
}

var _ Runtime = &dockerRuntime{}

func (d *dockerRuntime) Run(ctx context.Context, s *sequence.Step, opts ...RunOpt) error {
	ropts, err := createRunOpts(opts...)
	if err != nil {
		return err
	}

	pullOpts, err := d.getPullOpts(ctx, s, ropts)
	if err != nil {
		return err
	}

	contConf, err := d.getContainerConfig(ctx, s, ropts)
	if err != nil {
		return err
	}

	pullResp, err := d.client.ImagePull(ctx, contConf.Image, pullOpts)
	if err != nil {
		return err
	}
	defer pullResp.Close()
	// TODO make pretty
	io.Copy(os.Stdout, pullResp)

	workdir, mounts, err := d.getMounts(ctx, s, ropts)
	if err != nil {
		return err
	}

	contConf.WorkingDir = workdir
	hostConf, err := d.getHostConfig(ctx, s, ropts)
	if err != nil {
		return err
	}
	hostConf.Mounts = mounts

	createResp, err := d.client.ContainerCreate(ctx, contConf, hostConf, nil, nil, "")
	if err != nil {
		return err
	}

	attachResp, err := d.client.ContainerAttach(ctx, createResp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	// TODO make pretty, redirect stdout to json encoder if expecting response from step
	go stdcopy.StdCopy(os.Stdout, os.Stderr, attachResp.Reader)

	err = d.client.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusC, errC := d.client.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errC:
		if err != nil {
			return err
		}
	case <-statusC:
	case <-ctx.Done():
		defer d.client.ContainerRemove(ctx, createResp.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		return ctx.Err()
	}

	return nil
}
