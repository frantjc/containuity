package docker

import (
	"github.com/docker/docker/client"
	"github.com/frantjc/sequence/runtime"
)

type dockerVolume struct {
	name   string
	client *client.Client
}

var _ runtime.Volume = &dockerVolume{}

func (v *dockerVolume) Source() string {
	return v.name
}
