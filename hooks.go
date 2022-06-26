package sequence

type Hook[T any] func(T)

type Hooks[T any] []Hook[T]

func (h Hooks[T]) Hook(t T) {
	for _, hook := range h {
		hook(t)
	}
}
