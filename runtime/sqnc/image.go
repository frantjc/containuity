package sqnc

type sqncImage struct {
	ref string
}

func (i *sqncImage) Ref() string {
	return i.ref
}
