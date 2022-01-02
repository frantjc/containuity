package container

import (
	"context"

	"github.com/frantjc/sequence/pkg/sio"
)

type Container interface {
	Start(context.Context, *sio.Streams) error
	// Attach(context.Context) error
}
