package convert

import (
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
	"github.com/opencontainers/runtime-spec/specs-go"
	"google.golang.org/protobuf/types/known/anypb"
)

func MapInterfaceToAnyProto(i map[string]interface{}) map[string]*anypb.Any {
	a := map[string]*anypb.Any{}

	for k, v := range i {
		a[k] = v.(*anypb.Any)
	}

	return a
}

func RuntimeContainerToProtoContainer(container runtime.Container) *runtimev1.Container {
	return &runtimev1.Container{
		Id: container.GetID(),
	}
}

func RuntimeImageToProtoImage(image runtime.Image) *runtimev1.Image {
	return &runtimev1.Image{
		Ref: image.GetRef(),
	}
}

func RuntimeVolumeToProtoVolume(volume runtime.Volume) *runtimev1.Volume {
	return &runtimev1.Volume{
		Source: volume.GetSource(),
	}
}

func SpecsMountsToProtoMounts(m []specs.Mount) []*runtimev1.Mount {
	mounts := make([]*runtimev1.Mount, len(m))

	for i, j := range m {
		mounts[i] = &runtimev1.Mount{
			Source:      j.Source,
			Destination: j.Destination,
			Type:        j.Type,
			Options:     j.Options,
		}
	}

	return mounts
}
