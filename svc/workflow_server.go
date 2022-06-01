package svc

import (
	api "github.com/frantjc/sequence/api/v1/workflow"
	"github.com/frantjc/sequence/internal/conf"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/grpc"
)

func NewWorkflowService(runtime runtime.Runtime) (WorkflowService, error) {
	return &workflowServer{runtime: runtime}, nil
}

type workflowServer struct {
	api.UnimplementedWorkflowServer
	runtime runtime.Runtime
}

type WorkflowService interface {
	api.WorkflowServer
	Service
}

var _ WorkflowService = &workflowServer{}

func (s *workflowServer) RunWorkflow(in *api.RunWorkflowRequest, stream api.Workflow_RunWorkflowServer) error {
	var (
		ctx       = stream.Context()
		conf, err = conf.NewFromFlagsWithRepository(in.Repository)
		opts      = []workflow.ExecOpt{
			workflow.WithRuntime(s.runtime),
			workflow.WithGitHubToken(conf.GitHub.Token),
			workflow.WithRepository(in.Repository),
			workflow.WithStdout(grpcio.NewLogOutStreamWriter(stream)),
			workflow.WithStderr(grpcio.NewLogErrStreamWriter(stream)),
		}
	)
	if err != nil {
		return err
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
		return err
	}

	return executor.Start(ctx)
}

func (s *workflowServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterWorkflowServer(r, s)
}
