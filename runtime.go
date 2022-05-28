package sequence

import (
	"context"

	"github.com/frantjc/sequence/runtime"
)

func GetRuntime(names ...string) (runtime.Runtime, error) {
	return runtime.Get(context.Background(), names...)
}

func AnyRuntime() (runtime.Runtime, error) {
	return GetRuntime()
}
