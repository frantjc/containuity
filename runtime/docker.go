package runtime

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockermount "github.com/docker/docker/api/types/mount"
	dockervolume "github.com/docker/docker/api/types/volume"
	dockerclient "github.com/docker/docker/client"
	dockerstdcopy "github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/key"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func init() {
	sequence.RegisterRuntime("docker", func(ctx context.Context) (sequence.Runtime, error) {
		client, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}

		return &dockerRuntime{client}, nil
	})
}

type dockerRuntime struct {
	client *dockerclient.Client
}

var _ sequence.Runtime = &dockerRuntime{}

func (d *dockerRuntime) Run(ctx context.Context, s sequence.Steppable) error {
	steps, err := s.Steps()
	if err != nil {
		return err
	}
	log.Trace().Msgf("steppable has %d steps", len(steps))

	for _, step := range steps {
		ref, err := step.Image()
		if err != nil {
			return err
		}
		log.Trace().Msgf("image ref %s", ref)

		pullRes, err := d.client.ImagePull(ctx, ref, dockertypes.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer pullRes.Close()
		io.Copy(sequence.Stdout, pullRes)

		id := ctx.Value(key.Job).(string)
		if id == "" {
			id = strings.Replace(uuid.NewString(), "-", "", -1)
		}

		vol, err := d.client.VolumeCreate(ctx, dockervolume.VolumeCreateBody{
			Name: id,
		})
		if err != nil {
			return err
		}

		workdir := filepath.Join("/", "tmp", "step", id)
		createRes, err := d.client.ContainerCreate(ctx, &dockercontainer.Config{
			Image:        ref,
			Entrypoint:   step.Entrypoint,
			Cmd:          step.Cmd,
			WorkingDir:   workdir,
			AttachStdout: true,
			AttachStderr: true,
		}, &dockercontainer.HostConfig{
			AutoRemove: true,
			Mounts: []dockermount.Mount{
				{
					Source: vol.Name,
					Target: workdir,
					Type:   dockermount.TypeVolume,
				},
			},
			Privileged: step.Privileged,
		}, nil, nil, id)
		if err != nil {
			return err
		}

		attachRes, err := d.client.ContainerAttach(ctx, createRes.ID, dockertypes.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
			Logs:   true,
		})
		if err != nil {
			return err
		}
		defer attachRes.Close()
		go dockerstdcopy.StdCopy(sequence.Stdout, sequence.Stderr, attachRes.Reader)

		if err = d.client.ContainerStart(ctx, createRes.ID, dockertypes.ContainerStartOptions{}); err != nil {
			return err
		}

		statusC, errC := d.client.ContainerWait(ctx, createRes.ID, dockercontainer.WaitConditionNotRunning)
		select {
		case err := <-errC:
			if err != nil {
				return err
			}
		case status := <-statusC:
			if err := status.Error; err != nil {
				return fmt.Errorf(err.Message)
			} else if code := status.StatusCode; code != 0 {
				return fmt.Errorf("exit status %d", code)
			}
		case <-ctx.Done():
			defer d.client.ContainerKill(context.Background(), createRes.ID, "")
			return ctx.Err()
		}
	}

	return nil
}
