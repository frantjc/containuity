package docker

import "github.com/frantjc/sequence/runtime"

type dockerImage struct {
	ref string
}

var (
	_ runtime.Image = &dockerImage{}
)

func (i *dockerImage) Ref() string {
	return i.ref
}
