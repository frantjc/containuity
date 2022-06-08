package docker

import (
	"github.com/docker/docker/client"
	"github.com/frantjc/sequence/runtime"
)

type dockerContainer struct {
	id     string
	client *client.Client
}

var _ runtime.Container = &dockerContainer{}

func (c *dockerContainer) GetID() string {
	return c.id
}
