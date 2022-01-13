package runtime

import "context"

type Container interface {
	Exec(context.Context, ...ExecOpt) error
}
