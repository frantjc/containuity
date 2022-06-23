package runtime

import "context"

type Runtime interface {
	PullImage(context.Context, string) (Image, error)
	PruneImages(context.Context) error

	CreateContainer(context.Context, *Spec) (Container, error)
	GetContainer(context.Context, string) (Container, error)
	PruneContainers(context.Context) error

	CreateVolume(context.Context, string) (Volume, error)
	GetVolume(context.Context, string) (Volume, error)
	PruneVolumes(context.Context) error
}
