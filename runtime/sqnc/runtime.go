package sqnc

import (
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/v1/runtimev1connect"
)

type sqncRuntime struct {
	imageClient     runtimev1connect.ImageServiceClient
	containerClient runtimev1connect.ContainerServiceClient
	volumeClient    runtimev1connect.VolumeServiceClient
}

func NewRuntime(i runtimev1connect.ImageServiceClient, c runtimev1connect.ContainerServiceClient, v runtimev1connect.VolumeServiceClient) runtime.Runtime {
	return &sqncRuntime{i, c, v}
}

var _ runtime.Runtime = &sqncRuntime{}
