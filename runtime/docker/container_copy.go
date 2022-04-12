package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

func (c *dockerContainer) CopyTo(ctx context.Context, content io.Reader, destination string) error {
	return c.client.CopyToContainer(ctx, c.id, destination, content, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

func (c *dockerContainer) CopyFrom(ctx context.Context, source string) (io.ReadCloser, error) {
	content, _, err := c.client.CopyFromContainer(ctx, c.id, source)
	return content, err
}
