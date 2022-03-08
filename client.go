package sequence

import (
	"context"
	"io"

	jobapi "github.com/frantjc/sequence/api/v1/job"
	stepapi "github.com/frantjc/sequence/api/v1/step"
	workflowapi "github.com/frantjc/sequence/api/v1/workflow"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a wrapper around each of sequence's gRPC clients
type Client struct {
	jobclient      jobapi.JobClient
	stepclient     stepapi.StepClient
	workflowclient workflowapi.WorkflowClient
}

// New returns a new Client
func New(ctx context.Context, addr string, opts ...ClientOpt) (*Client, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		jobclient:      jobapi.NewJobClient(cc),
		stepclient:     stepapi.NewStepClient(cc),
		workflowclient: workflowapi.NewWorkflowClient(cc),
	}, nil
}

// JobClient returns the client's underlying gRPC JobClient
func (c *Client) JobClient() jobapi.JobClient {
	return c.jobclient
}

// StepClient returns the client's underlying gRPC StepClient
func (c *Client) StepClient() stepapi.StepClient {
	return c.stepclient
}

// WorkflowClient returns the client's underlying gRPC WorkflowClient
func (c *Client) WorkflowClient() workflowapi.WorkflowClient {
	return c.workflowclient
}

// RunStep calls the underlying gRPC StepClient's RunStep and
// writes its logs to the given io.Writer
func (c *Client) RunStep(ctx context.Context, step *workflow.Step, w io.Writer) error {
	stream, err := c.StepClient().RunStep(ctx, &stepapi.RunStepRequest{
		Step: convert.StepToTypeStep(step),
	})
	if err != nil {
		return err
	}

	for {
		l, err := stream.Recv()
		if err == io.EOF {
			return stream.CloseSend()
		} else if err != nil {
			return err
		}

		w.Write([]byte(l.Line))
	}
}

func (c *Client) RunJob(ctx context.Context, job *workflow.Job, w io.Writer) error {
	stream, err := c.JobClient().RunJob(ctx, &jobapi.RunJobRequest{
		Job: convert.JobToTypeJob(job),
	})
	if err != nil {
		return err
	}

	for {
		l, err := stream.Recv()
		if err == io.EOF {
			return stream.CloseSend()
		} else if err != nil {
			return err
		}

		w.Write([]byte(l.Line))
	}
}

func (c *Client) RunWorkflow(ctx context.Context, workflow *workflow.Workflow, w io.Writer) error {
	stream, err := c.WorkflowClient().RunWorkflow(ctx, &workflowapi.RunWorkflowRequest{
		Workflow: convert.WorkflowToTypeWorkflow(workflow),
	})
	if err != nil {
		return err
	}

	for {
		l, err := stream.Recv()
		if err == io.EOF {
			return stream.CloseSend()
		} else if err != nil {
			return err
		}

		w.Write([]byte(l.Line))
	}
}
