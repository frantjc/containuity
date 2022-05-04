package runtime

import (
	"context"
	"io"
)

type Container interface {
	ID() string
	Exec(context.Context, *Streams) error
	CopyTo(context.Context, io.Reader, string) error
	CopyFrom(context.Context, string) (io.ReadCloser, error)
	Remove(context.Context) error
	Start(context.Context) error
	Attach(context.Context, *Streams) error
}
