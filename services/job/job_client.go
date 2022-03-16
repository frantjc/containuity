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
		conf, _ = conf.Get()
		stream  = grpcio.NewLogStream(ctx)
		opts    = []workflow.RunOpt{
			workflow.WithStdout(grpcio.NewLogStreamWriter(stream)),
			workflow.WithGitHubToken(conf.GitHub.Token),
			workflow.WithRunnerImage(conf.Runtime.Image),
			workflow.WithWorkdir(conf.RootDir),
		}
	)

	if conf.Verbose {
		opts = append(opts, workflow.WithVerbose)
	}

	go func() {
		defer stream.CloseSend()
		if err := workflow.RunJob(ctx, c.runtime, convert.ProtoJobToJob(in.Job), opts...); err != nil {
			stream.SendErr(err)
		}
	}()

	return stream, nil
}
