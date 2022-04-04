package job

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/job"
	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/grpc"
)

type jobClient struct {
	runtime runtime.Runtime
}

var _ api.JobClient = &jobClient{}

func (c *jobClient) RunJob(ctx context.Context, in *api.RunJobRequest, _ ...grpc.CallOption) (api.Job_RunJobClient, error) {
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

	if in.Job != nil {
		opts = append(opts, workflow.WithJob(
			convert.ProtoJobToJob(in.Job),
		))
	}

	executor, err := workflow.NewJobExecutor(convert.ProtoJobToJob(in.Job), opts...)
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
