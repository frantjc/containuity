package sqnc

type sqncVolume struct {
	source string
	client RuntimeServiceClient
}

func (v *sqncVolume) GetSource() string {
	return v.source
}
