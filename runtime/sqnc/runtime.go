package sqnc

import (
	containerapi "github.com/frantjc/sequence/api/v1/container"
	imageapi "github.com/frantjc/sequence/api/v1/image"
	volumeapi "github.com/frantjc/sequence/api/v1/volume"

	"github.com/frantjc/sequence/runtime"
)

type sqncRuntime struct {
	imageClient     imageapi.ImageClient
	containerClient containerapi.ContainerClient
	volumeClient    volumeapi.VolumeClient
}

func NewRuntime(i imageapi.ImageClient, c containerapi.ContainerClient, v volumeapi.VolumeClient) runtime.Runtime {
	return &sqncRuntime{i, c, v}
}

var _ runtime.Runtime = &sqncRuntime{}
