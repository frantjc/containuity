package sequence

import (
	"context"
	"io"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/frantjc/sequence/internal/protobufio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/sqnc"
	"github.com/frantjc/sequence/runtime/v1/runtimev1connect"
	workflowv1 "github.com/frantjc/sequence/workflow/v1"
	"github.com/frantjc/sequence/workflow/v1/workflowv1connect"
)

// Client is a wrapper around each of sequence's rpc clients
type Client struct {
	httpClient       *http.Client
	workflowsClient  workflowv1connect.WorkflowServiceClient
	containersClient runtimev1connect.ContainerServiceClient
	imagesClient     runtimev1connect.ImageServiceClient
	volumesClient    runtimev1connect.VolumeServiceClient
}

// New is an alias to NewClient
var New = NewClient

// NewClient returns a new Client
func NewClient(ctx context.Context, addr string, opts ...ClientOpt) (*Client, error) {
	client := &Client{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	client.workflowsClient = workflowv1connect.NewWorkflowServiceClient(client.httpClient, addr)
	client.containersClient = runtimev1connect.NewContainerServiceClient(client.httpClient, addr)
	client.imagesClient = runtimev1connect.NewImageServiceClient(client.httpClient, addr)
	client.volumesClient = runtimev1connect.NewVolumeServiceClient(client.httpClient, addr)

	return client, nil
}

// WorkflowsClient returns the client's underlying rpc WorkflowClient
func (c *Client) WorkflowsClient() workflowv1connect.WorkflowServiceClient {
	return c.workflowsClient
}

// ContainerClient returns the client's underlying rpc ContainerClient
func (c *Client) ContainersClient() runtimev1connect.ContainerServiceClient {
	return c.containersClient
}

// ImageClient returns the client's underlying rpc ImageClient
func (c *Client) ImagesClient() runtimev1connect.ImageServiceClient {
	return c.imagesClient
}

// VolumeClient returns the client's underlying rpc VolumeClient
func (c *Client) VolumesClient() runtimev1connect.VolumeServiceClient {
	return c.volumesClient
}

// Runtime returns a runtime.Runtime implementation using the underlying clients
func (c *Client) Runtime() runtime.Runtime {
	return sqnc.NewRuntime(c.ImagesClient(), c.ContainersClient(), c.VolumesClient())
}

// RunStep calls the underlying rpc StepClient's RunStep and
// writes its logs to the given io.Writer
func (c *Client) RunStep(ctx context.Context, step *workflowv1.Step, w io.Writer, opts ...RunOpt) error {
	ro := defaultRunOpts()
	for _, opt := range opts {
		err := opt(ro)
		if err != nil {
			return err
		}
	}

	stream, err := c.WorkflowsClient().RunStep(ctx, connect.NewRequest(&workflowv1.RunStepRequest{
		Step:        step,
		Job:         ro.job,
		Workflow:    ro.workflow,
		Repository:  ro.repository,
		RunnerImage: ro.runnerImage,
		Verbose:     ro.verbose,
	}))
	if err != nil {
		return err
	}

	return protobufio.DemultiplexLogStream[*workflowv1.RunStepResponse](stream, w, w)
}

// RunJob calls the underlying rpc JobClient's RunJob and
// writes its logs to the given io.Writer
func (c *Client) RunJob(ctx context.Context, job *workflowv1.Job, w io.Writer, opts ...RunOpt) error {
	ro := defaultRunOpts()
	for _, opt := range opts {
		err := opt(ro)
		if err != nil {
			return err
		}
	}

	stream, err := c.WorkflowsClient().RunJob(ctx, connect.NewRequest(&workflowv1.RunJobRequest{
		Job:         job,
		Workflow:    ro.workflow,
		Repository:  ro.repository,
		RunnerImage: ro.runnerImage,
		Verbose:     ro.verbose,
	}))
	if err != nil {
		return err
	}

	return protobufio.DemultiplexLogStream[*workflowv1.RunJobResponse](stream, w, w)
}

// RunWorkflow calls the underlying rpc WorkflowClient's RunWorkflow and
// writes its logs to the given io.Writer
func (c *Client) RunWorkflow(ctx context.Context, workflow *workflowv1.Workflow, w io.Writer, opts ...RunOpt) error {
	ro := defaultRunOpts()
	for _, opt := range opts {
		err := opt(ro)
		if err != nil {
			return err
		}
	}

	stream, err := c.WorkflowsClient().RunWorkflow(ctx, connect.NewRequest(&workflowv1.RunWorkflowRequest{
		Workflow:    workflow,
		Repository:  ro.repository,
		RunnerImage: ro.runnerImage,
		Verbose:     ro.verbose,
	}))
	if err != nil {
		return err
	}

	return protobufio.DemultiplexLogStream[*workflowv1.RunWorkflowResponse](stream, w, w)
}
