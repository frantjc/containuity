package sequence

import "net/http"

type ClientOpt func(*Client) error

func WithHTTPClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}
