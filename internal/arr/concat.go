package arr

func Concat[T any](in ...[]T) []T {
	out := []T{}
	for _, t := range in {
		out = append(out, t...)
	}

	return out
}
