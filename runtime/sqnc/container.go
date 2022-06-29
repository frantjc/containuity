package sqnc

type sqncContainer struct {
	id     string
	client RuntimeServiceClient
}

func (c *sqncContainer) GetID() string {
	return c.id
}
