package docker

import "github.com/docker/docker/client"

type dockerVolume struct {
	name   string
	client *client.Client
}

func (v *dockerVolume) GetSource() string {
	return v.name
}
