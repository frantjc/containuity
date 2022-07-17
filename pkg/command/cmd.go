package command

import "context"

type Cmd interface {
	Execute() error
	ExecuteContext(context.Context) error
}
