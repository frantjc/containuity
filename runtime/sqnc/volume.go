package sqnc

import volumeapi "github.com/frantjc/sequence/pb/v1/volume"

type sqncVolume struct {
	source string
	client volumeapi.VolumeClient
}

func (v *sqncVolume) Source() string {
	return v.source
}
