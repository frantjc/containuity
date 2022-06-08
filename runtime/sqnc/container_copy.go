package sqnc

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/bufbuild/connect-go"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
)

func (c *sqncContainer) CopyTo(ctx context.Context, content io.Reader, destination string) error {
	b, err := ioutil.ReadAll(content)
	if err != nil {
		return err
	}

	_, err = c.client.CopyToContainer(ctx, connect.NewRequest(&runtimev1.CopyToContainerRequest{
		Id:          c.id,
		Content:     b,
		Destination: destination,
	}))

	return err
}

func (c *sqncContainer) CopyFrom(ctx context.Context, source string) (io.ReadCloser, error) {
	res, err := c.client.CopyFromContainer(ctx, connect.NewRequest(&runtimev1.CopyFromContainerRequest{
		Id: c.id,
	}))

	return io.NopCloser(
		bytes.NewReader(res.Msg.GetContent()),
	), err
}
