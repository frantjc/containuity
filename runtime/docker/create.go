package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/frantjc/sequence/defaults"
	"github.com/frantjc/sequence/runtime"
	"github.com/google/go-containerregistry/pkg/name"
)

func (r *dockerRuntime) Create(ctx context.Context, opts ...runtime.SpecOpt) (runtime.Container, error) {
	spec, err := runtime.NewSpec(opts...)
	if err != nil {
		return nil, err
	}

	pref, err := name.ParseReference(spec.Image)
	if err != nil {
		return nil, err
	}

	conf := &container.Config{
		Image:      pref.Name(),
		Entrypoint: spec.Entrypoint,
		Cmd:        spec.Cmd,
		WorkingDir: spec.Cwd,
		Env:        spec.Env,
		Labels:     defaults.Labels,
	}

	hconf := &container.HostConfig{
		AutoRemove: true,
		Privileged: spec.Privileged,
	}

	for _, m := range spec.Mounts {
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
