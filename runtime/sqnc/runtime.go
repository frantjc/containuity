package sqnc

import (
	"github.com/frantjc/sequence/runtime"
)

const RuntimeName = "sqnc"

func NewRuntime(c RuntimeServiceClient) runtime.Runtime {
	return &sqncRuntime{c}
}

type sqncRuntime struct {
	runtimeClient RuntimeServiceClient
}
