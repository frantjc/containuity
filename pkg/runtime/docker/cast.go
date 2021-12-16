package runtime

import (
	"github.com/frantjc/sequence/pkg/container"
	"github.com/frantjc/sequence/pkg/runtime"
)

var (
	_ runtime.Runtime     = &dockerRuntime{}
	_ container.Container = &dockerContainer{}
)
