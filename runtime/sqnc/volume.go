package sqnc

import "github.com/frantjc/sequence/runtime/v1/runtimev1connect"

type sqncVolume struct {
	source string
	client runtimev1connect.VolumeServiceClient
}

func (v *sqncVolume) GetSource() string {
	return v.source
}
