package docker

import "github.com/docker/docker/client"

type dockerContainer struct {
	id     string
	client *client.Client
}

func (c *dockerContainer) GetID() string {
	return c.id
}
