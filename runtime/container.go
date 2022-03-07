package runtime

import "context"

type Container interface {
	ID() string
	Exec(context.Context, ...ExecOpt) error
}
