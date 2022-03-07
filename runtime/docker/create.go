package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/frantjc/sequence/defaults"
	"github.com/frantjc/sequence/runtime"
	"github.com/google/go-containerregistry/pkg/name"
)

func (r *dockerRuntime) CreateContainer(ctx context.Context, s *runtime.Spec) (runtime.Container, error) {
	pref, err := name.ParseReference(s.Image)
	if err != nil {
		return nil, err
	}

	conf := &container.Config{
		Image:      pref.Name(),
		Entrypoint: s.Entrypoint,
		Cmd:        s.Cmd,
		WorkingDir: s.Cwd,
		Env:        s.Env,
		Labels:     defaults.Labels,
	}

	hconf := &container.HostConfig{
		AutoRemove: true,
		Privileged: s.Privileged,
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
				Labels: defaults.Labels,
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
