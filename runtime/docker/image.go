package docker

import "github.com/frantjc/sequence/runtime"

type dockerImage struct {
	ref string
}

var (
	_ runtime.Image = &dockerImage{}
)

func (i *dockerImage) GetRef() string {
	return i.ref
}
