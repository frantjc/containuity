package sqnc

import containerapi "github.com/frantjc/sequence/api/v1/container"

type sqncContainer struct {
	id     string
	client containerapi.ContainerClient
}

func (c *sqncContainer) ID() string {
	return c.id
}
