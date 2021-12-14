package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/frantjc/sequence"
)

var (
	defaultDockerMounts = []mount.Mount{
		{
			Type:   mount.TypeTmpfs,
			Target: "/dev/shm",
			TmpfsOptions: &mount.TmpfsOptions{
				Mode: 01777,
			},
		},
		{
			Type:   mount.TypeBind,
			Source: filepath.Join(os.TempDir()),
			Target: filepath.Join("/tmp"),
		},
	}
)

func (d *dockerRuntime) getMounts(ctx context.Context, s *sequence.Step, r *runOpts) (string, []mount.Mount, error) {
	var (
		mounts  = defaultDockerMounts
		id      = r.id
		workdir = filepath.Join("/tmp", id, "workdir")
	)

	vol, err := d.client.VolumeCreate(ctx, volume.VolumeCreateBody{
		Name:   fmt.Sprintf("%s-workdir", id),
		Labels: labels,
	})
	if err != nil {
		return workdir, mounts, err
	}

	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeVolume,
		Source: vol.Name,
		Target: workdir,
	})

	if s.Uses != "" {
		vol, err := d.client.VolumeCreate(ctx, volume.VolumeCreateBody{
			Name:   fmt.Sprintf("%s-action", id),
			Labels: labels,
		})
		if err != nil {
			return workdir, mounts, err
		}

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: vol.Name,
			Target: filepath.Join("/tmp", "action"),
		})

		vol, err = d.client.VolumeCreate(ctx, volume.VolumeCreateBody{
			Name:   fmt.Sprintf("%s-tool-cache", id),
			Labels: labels,
		})
		if err != nil {
			return workdir, mounts, err
		}

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: vol.Name,
			Target: filepath.Join("/tmp", "tool-cache"),
		})
	}

	return workdir, mounts, nil
}
