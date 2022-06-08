package sqnc

import "github.com/frantjc/sequence/runtime/v1/runtimev1connect"

type sqncContainer struct {
	id     string
	client runtimev1connect.ContainerServiceClient
}

func (c *sqncContainer) GetID() string {
	return c.id
}
