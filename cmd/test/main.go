package main

import (
	"context"
	"io"
	"os"

	dtypes "github.com/docker/docker/api/types"
	dcontainer "github.com/docker/docker/api/types/container"
	dclient "github.com/docker/docker/client"
	dstdcopy "github.com/docker/docker/pkg/stdcopy"
)

func main() {
	var (
		ctx         = context.Background()
		client, err = dclient.NewClientWithOpts(dclient.FromEnv, dclient.WithAPIVersionNegotiation())
		image       = "docker.io/library/alpine"
	)
	if err != nil {
		panic(err)
	}

	pullResp, err := client.ImagePull(ctx, image, dtypes.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, pullResp)
	pullResp.Close()

	createResp, err := client.ContainerCreate(ctx, &dcontainer.Config{
		Image: image,
		Cmd:   []string{"whoami"},
		Env:   []string{"TEST=test"},
		User:  "65534:65534",
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	attachResp, err := client.ContainerAttach(ctx, createResp.ID, dtypes.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		panic(err)
	}
	// TODO should we do this?
	// go io.Copy(attachResp.Conn, c.opts.Stdin)
	// TODO prettify this
	go dstdcopy.StdCopy(os.Stdout, os.Stderr, attachResp.Reader)
	defer attachResp.Close()

	err = client.ContainerStart(ctx, createResp.ID, dtypes.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	statusC, errC := client.ContainerWait(ctx, createResp.ID, dcontainer.WaitConditionNotRunning)
	select {
	case err := <-errC:
		if err != nil {
			panic(err)
		}
	case <-ctx.Done():
		defer client.ContainerRemove(ctx, createResp.ID, dtypes.ContainerRemoveOptions{
			Force: true,
		})
		panic(err)
	case <-statusC:
	}

	os.Exit(0)
}
