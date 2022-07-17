package sequence

import "github.com/frantjc/sequence/pkg/github/actions"

type Event[T any] struct {
	Type          T
	GlobalContext *actions.GlobalContext
}

type Hook[T any] func(*Event[T])

type Hooks[T any] []Hook[T]

func (h Hooks[T]) Invoke(event *Event[T]) {
	for _, hook := range h {
		hook(event)
	}
}
