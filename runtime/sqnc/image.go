package sqnc

type sqncImage struct {
	ref string
}

func (i *sqncImage) GetRef() string {
	return i.ref
}
