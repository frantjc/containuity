package sequence

import (
	"context"
	"net/http"
)

type Opt func(context.Context, *Client) error

func WithHTTPClient(httpClient *http.Client) Opt {
	return func(ctx context.Context, c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}
