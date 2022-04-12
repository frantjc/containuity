package runtime

import "context"

type Volume interface {
	Source() string
	Remove(context.Context) error
}
