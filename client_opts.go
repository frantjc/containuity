package sequence

import (
	"context"
	"net/http"
)

type ClientOpt func(context.Context, *Client) error

func WithHTTPClient(hc *http.Client) ClientOpt {
	return func(ctx context.Context, c *Client) error {
		c.httpClient = hc
		return nil
	}
}
