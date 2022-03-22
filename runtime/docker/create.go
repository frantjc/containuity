package docker

import (
	"context"
	"fmt"
	"os/exec"
	goruntime "runtime"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/frantjc/sequence/runtime"
)

func (r *dockerRuntime) CreateContainer(ctx context.Context, s *runtime.Spec) (runtime.Container, error) {
	pref, err := reference.Parse(s.Image)
	if err != nil {
		return nil, err
	}

	var (
		addr = r.client.DaemonHost()
		conf = &container.Config{
			Image:      pref.String(),
			Entrypoint: s.Entrypoint,
			Cmd:        s.Cmd,
			WorkingDir: s.Cwd,
			Env:        append(s.Env, fmt.Sprintf("DOCKER_HOST=%s", addr)),
			Labels:     labels,
		}
		hconf = &container.HostConfig{
			AutoRemove: true,
			Privileged: s.Privileged,
		}
	)

	if strings.HasPrefix(addr, "unix://") {
		sock := strings.TrimPrefix(addr, "unix://")
		hconf.Mounts = append(hconf.Mounts, mount.Mount{
			Source: sock,
			Target: "/var/run/docker.sock",
			Type:   runtime.MountTypeBind,
		})
	}

	if goruntime.GOOS == "linux" {
		if docker, err := exec.LookPath("docker"); err == nil {
			hconf.Mounts = append(hconf.Mounts, mount.Mount{
				Source: docker,
				Target: "/usr/local/bin/docker",
				Type:   runtime.MountTypeBind,
			})
		}
	}

	for _, m := range s.Mounts {
		dt := mount.Type(m.Type)
		dm := mount.Mount{
			Type:   dt,
			Source: m.Source,
			Target: m.Destination,
		}

		switch dt {
		case mount.TypeBind:
		case mount.TypeVolume:
			dm.VolumeOptions = &mount.VolumeOptions{
				Labels: labels,
			}
		case mount.TypeTmpfs:
			dm.Source = ""
		}

		for _, opt := range m.Options {
			switch opt {
			case runtime.MountOptReadOnly:
				dm.ReadOnly = true
			}
		}

		hconf.Mounts = append(hconf.Mounts, dm)
	}

	createResp, err := r.client.ContainerCreate(ctx, conf, hconf, nil, nil, "")
	if err != nil {
		return nil, err
	}

	return &dockerContainer{id: createResp.ID, client: r.client}, nil
}
