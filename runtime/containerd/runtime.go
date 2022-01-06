package containerd

import (
	"github.com/containerd/containerd"
)

type containerdRuntime struct {
	client *containerd.Client
}
