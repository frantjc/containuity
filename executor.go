package sequence

import (
	"context"
	"io"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/shim"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/runtimeutil"
)

type Executor struct {
	RunnerImage       runtime.Image
	Runtime           runtime.Runtime
	Stdout, Stderr    io.Writer
	Stdin             io.Reader
	Verbose           bool
	GlobalContext     *actions.GlobalContext
	OnImagePull       Hooks[runtime.Image]
	OnContainerCreate Hooks[runtime.Container]
	OnVolumeCreate    Hooks[runtime.Volume]
	OnWorkflowCommand Hooks[*actions.WorkflowCommand]
}

func (e Executor) RunContainer(ctx context.Context, spec *runtime.Spec, streams *runtime.Streams) error {
	image, err := e.Runtime.PullImage(ctx, spec.Image)
	if err != nil {
		return err
	}
	e.OnImagePull.Hook(image)

	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeVolume {
			volume, err := e.Runtime.CreateVolume(ctx, mount.Source)
			if err != nil {
				if volume, err = e.Runtime.GetVolume(ctx, mount.Source); err != nil {
					return err
				}
			} else {
				e.OnVolumeCreate.Hook(volume)
			}
			mount.Source = volume.GetSource()
		}
	}

	container, err := e.Runtime.CreateContainer(ctx, spec)
	if err != nil {
		return err
	}
	e.OnContainerCreate.Hook(container)

	tarArchive, err := runtimeutil.NewSingleFileTarArchiveReader(shimName, shim.Bytes)
	if err != nil {
		return err
	}

	if err = container.CopyTo(ctx, tarArchive, shimDir); err != nil {
		return err
	}

	return container.Exec(ctx, streams)
}