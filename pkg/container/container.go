package container

import (
	"context"

	"github.com/frantjc/sequence/pkg/io"
)

type Container interface {
	Start(context.Context, io.Streams) error
	// Attach(context.Context) error
}
