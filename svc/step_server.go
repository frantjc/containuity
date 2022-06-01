package svc

import (
	api "github.com/frantjc/sequence/api/v1/step"
	"github.com/frantjc/sequence/internal/conf"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/grpc"
)

func NewStepService(runtime runtime.Runtime) (StepService, error) {
	return &stepServer{runtime: runtime}, nil
}

type stepServer struct {
	api.UnimplementedStepServer
	runtime runtime.Runtime
}

type StepService interface {
	api.StepServer
	Service
}

var _ StepService = &stepServer{}

func (s *stepServer) RunStep(in *api.RunStepRequest, stream api.Step_RunStepServer) error {
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

	if in.Job != nil {
		opts = append(opts, workflow.WithJob(
			convert.ProtoJobToJob(in.Job),
		))
	}

	executor, err := workflow.NewStepExecutor(convert.ProtoStepToStep(in.Step), opts...)
	if err != nil {
		return err
	}

	return executor.Start(ctx)
}

func (s *stepServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterStepServer(r, s)
}
