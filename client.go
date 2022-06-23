package sequence

import (
	"context"
	"net/http"
	"net/url"

	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/sqnc"
)

// Client is a wrapper around each of sequence's rpc clients
type Client struct {
	httpClient      *http.Client
	workflowsClient WorkflowServiceClient
	runtimeClient   sqnc.RuntimeServiceClient
}

// New is an alias to NewClient
var New = NewClient

// NewClient returns a new Client
func NewClient(ctx context.Context, addr *url.URL, opts ...Opt) (*Client, error) {
	client := &Client{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		if err := opt(ctx, client); err != nil {
			return nil, err
		}
	}

	client.workflowsClient = NewWorkflowServiceClient(client.httpClient, addr.String())
	client.runtimeClient = sqnc.NewRuntimeServiceClient(client.httpClient, addr.String())

	return client, nil
}

// WorkflowsClient returns the client's underlying rpc WorkflowClient
func (c *Client) WorkflowsClient() WorkflowServiceClient {
	return c.workflowsClient
}

// RuntimeClient returns the client's underlying rpc RuntimeClient
func (c *Client) RuntimeClient() sqnc.RuntimeServiceClient {
	return c.runtimeClient
}

// Runtime returns a runtime.Runtime implementation using the underlying clients
func (c *Client) Runtime() runtime.Runtime {
	return sqnc.NewRuntime(c.RuntimeClient())
}
