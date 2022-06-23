package convert

import (
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/sqnc"
	"google.golang.org/protobuf/types/known/anypb"
)

func MapInterfaceToAnyProto(i map[string]interface{}) map[string]*anypb.Any {
	a := map[string]*anypb.Any{}

	for k, v := range i {
		a[k] = v.(*anypb.Any)
	}

	return a
}

func RuntimeContainerToProtoContainer(container runtime.Container) *sqnc.Container {
	return &sqnc.Container{
		Id: container.GetID(),
	}
}

func RuntimeImageToProtoImage(image runtime.Image) *sqnc.Image {
	return &sqnc.Image{
		Ref: image.GetRef(),
	}
}

func RuntimeVolumeToProtoVolume(volume runtime.Volume) *sqnc.Volume {
	return &sqnc.Volume{
		Source: volume.GetSource(),
	}
}
