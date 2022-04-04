package workflow

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/workflow"
	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/grpc"
)

type workflowClient struct {
	runtime runtime.Runtime
}

var _ api.WorkflowClient = &workflowClient{}

func (c *workflowClient) RunWorkflow(ctx context.Context, in *api.RunWorkflowRequest, _ ...grpc.CallOption) (api.Workflow_RunWorkflowClient, error) {
	var (
		conf, err = conf.NewFromFlagsWithRepository(in.Repository)
		stream    = grpcio.NewLogStream(ctx)
		opts      = []workflow.ExecOpt{
			workflow.WithRuntime(c.runtime),
			workflow.WithGitHubToken(conf.GitHub.Token),
			workflow.WithRepository(in.Repository),
			workflow.WithWorkdir(conf.RootDir),
			workflow.WithStdout(grpcio.NewLogOutStreamWriter(stream)),
			workflow.WithStderr(grpcio.NewLogErrStreamWriter(stream)),
		}
	)
	if err != nil {
		return nil, err
	}

	if in.Verbose || conf.Verbose {
		opts = append(opts, workflow.WithVerbose)
	}

	if in.RunnerImage != "" {
		opts = append(opts, workflow.WithRunnerImage(in.RunnerImage))
	} else {
		in.RunnerImage = conf.Runtime.RunnerImage
	}

	executor, err := workflow.NewWorkflowExecutor(convert.ProtoWorkflowToWorkflow(in.Workflow), opts...)
	if err != nil {
		return nil, err
	}

	go func() {
		defer stream.CloseSend()
		if err = executor.Start(ctx); err != nil {
			stream.SendErr(err)
		}
	}()

	return stream, nil
}
