package runtime

import "context"

type Volume interface {
	GetSource() string
	Remove(context.Context) error
}
