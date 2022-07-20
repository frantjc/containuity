package sequence

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/internal/shim"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/runtimeutil"
)

type Executor interface {
	Execute() error
	ExecuteContext(context.Context) error
}

type executor struct {
	RunnerImage          runtime.Image
	Runtime              runtime.Runtime
	StreamOut, StreamErr io.Writer
	StreamIn             io.Reader
	Verbose              bool
	GlobalContext        *actions.GlobalContext
	OnImagePull          Hooks[runtime.Image]
	OnContainerCreate    Hooks[runtime.Container]
	OnVolumeCreate       Hooks[runtime.Volume]
	OnWorkflowCommand    Hooks[*actions.WorkflowCommand]
	OnStepStart          Hooks[*Step]
	OnStepFinish         Hooks[*Step]
	OnJobStart           Hooks[*Job]
	OnJobFinish          Hooks[*Job]
	OnWorkflowStart      Hooks[*Workflow]
	OnWorkflowFinish     Hooks[*Workflow]
}

func (e *executor) GetID() string {
	id := e.GlobalContext.GitHubContext.Job

	if e.GlobalContext.GitHubContext.Workflow != "" {
		id = fmt.Sprintf("%s-%s", e.GlobalContext.GitHubContext.Workflow, id)
	}

	return strings.ToLower(id)
}

func (e *executor) RunContainer(ctx context.Context, spec *runtime.Spec, streams *runtime.Streams) error {
	image, err := e.Runtime.PullImage(ctx, spec.Image)
	if err != nil {
		return err
	}

	e.OnImagePull.Invoke(&Event[runtime.Image]{
		Type:          image,
		GlobalContext: e.GlobalContext,
	})

	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeVolume {
			volume, err := e.Runtime.CreateVolume(ctx, mount.Source)
			if err != nil {
				if volume, err = e.Runtime.GetVolume(ctx, mount.Source); err != nil {
					return err
				}
			} else {
				e.OnVolumeCreate.Invoke(&Event[runtime.Volume]{
					Type:          volume,
					GlobalContext: e.GlobalContext,
				})
			}
			mount.Source = volume.GetSource()
		}
	}

	container, err := e.Runtime.CreateContainer(ctx, spec)
	if err != nil {
		return err
	}

	tarArchive, err := runtimeutil.NewSingleFileTarArchiveReader(shim.Name, shim.Bytes)
	if err != nil {
		return err
	}

	if err = container.CopyTo(ctx, tarArchive, path.Dir(paths.Shim)); err != nil {
		return err
	}

	e.OnContainerCreate.Invoke(&Event[runtime.Container]{
		Type:          container,
		GlobalContext: e.GlobalContext,
	})

	return container.Exec(ctx, streams)
}
