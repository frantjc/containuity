package runtime

import (
	"context"
	"io"
)

type Container interface {
	GetID() string
	Exec(context.Context, *Exec, *Streams) error
	CopyTo(context.Context, io.Reader, string) error
	CopyFrom(context.Context, string) (io.ReadCloser, error)
	Stop(context.Context) error
	Remove(context.Context) error
	Start(context.Context) error
	Attach(context.Context, *Streams) error
}
