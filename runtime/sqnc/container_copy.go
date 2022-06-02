package sqnc

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	containerapi "github.com/frantjc/sequence/pb/v1/container"
)

func (c *sqncContainer) CopyTo(ctx context.Context, content io.Reader, destination string) error {
	b, err := ioutil.ReadAll(content)
	if err != nil {
		return err
	}

	_, err = c.client.CopyToContainer(ctx, &containerapi.CopyToContainerRequest{
		Id:          c.id,
		Content:     b,
		Destination: destination,
	})
	return err
}

func (c *sqncContainer) CopyFrom(ctx context.Context, source string) (io.ReadCloser, error) {
	res, err := c.client.CopyFromContainer(ctx, &containerapi.CopyFromContainerRequest{
		Id: c.id,
	})
	return io.NopCloser(
		bytes.NewReader(res.Content),
	), err
}
