package services

import (
	api "github.com/frantjc/sequence/api/v1/job"
	"github.com/frantjc/sequence/internal/conf"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/grpc"
)

func NewJobService(runtime runtime.Runtime) (JobService, error) {
	svc := &jobServer{
		svc: &service{runtime},
	}

	return svc, nil
}

type jobServer struct {
	api.UnimplementedJobServer
	svc *service
}

type JobService interface {
	api.JobServer
	Service
}

var _ JobService = &jobServer{}

func (s *jobServer) RunJob(in *api.RunJobRequest, stream api.Job_RunJobServer) error {
	var (
		ctx       = stream.Context()
		conf, err = conf.NewFromFlagsWithRepository(in.Repository)
		opts      = []workflow.ExecOpt{
			workflow.WithRuntime(s.svc.runtime),
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

	executor, err := workflow.NewJobExecutor(convert.ProtoJobToJob(in.Job), opts...)
	if err != nil {
		return err
	}

	return executor.Start(ctx)
}

func (s *jobServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterJobServer(r, s)
}
