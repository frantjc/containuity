package runtimes

import "errors"

var ErrRuntimeNotFound = errors.New("runtime not found")

func ErrIsRuntimeNotFound(err error) bool {
	return errors.Is(err, ErrRuntimeNotFound)
}
