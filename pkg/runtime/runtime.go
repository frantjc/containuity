package runtime

import (
	"context"

	"github.com/frantjc/sequence/pkg/container"
)

type Runtime interface {
	Pull(context.Context, string) error
	Create(context.Context, *container.Spec) (container.Container, error)
}
