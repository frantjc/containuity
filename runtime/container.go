package runtime

import (
	"context"
	"io"
)

type Container interface {
	ID() string
	Exec(context.Context, *Exec) error
	CopyTo(context.Context, io.Reader, string) error
	CopyFrom(context.Context, string) (io.ReadCloser, error)
}
