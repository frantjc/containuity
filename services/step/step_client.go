package step

import (
	"context"

	api "github.com/frantjc/sequence/api/v1/step"
	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/grpc"
)

type stepClient struct {
	runtime runtime.Runtime
}

var _ api.StepClient = &stepClient{}

func (c *stepClient) RunStep(ctx context.Context, in *api.RunStepRequest, _ ...grpc.CallOption) (api.Step_RunStepClient, error) {
	var (
		conf, err = conf.NewFromFlagsWithRepository(in.Repository)
		stream    = grpcio.NewLogStream(ctx)
		opts      = []workflow.RunOpt{
			workflow.WithStdout(grpcio.NewLogOutStreamWriter(stream)),
			workflow.WithGitHubToken(conf.GitHub.Token),
			workflow.WithWorkdir(conf.RootDir),
			workflow.WithSecrets(conf.Secrets),
		}
	)
	if err != nil {
		return nil, err
	}

	if in.Repository != "" {
		opts = append(opts, workflow.WithRepository(in.Repository))
	}

	if in.RunnerImage != "" {
		opts = append(opts, workflow.WithRunnerImage(in.RunnerImage))
	} else {
		opts = append(opts, workflow.WithRunnerImage(conf.Runtime.RunnerImage))
	}

	if in.ActionImage != "" {
		opts = append(opts, workflow.WithActionImage(in.ActionImage))
	}

	if conf.Verbose || in.Verbose {
		opts = append(opts, workflow.WithVerbose)
	}

	go func() {
		defer stream.CloseSend()
		if ctx, _, err = workflow.RunStep(ctx, c.runtime, convert.ProtoStepToStep(in.Step), opts...); err != nil {
			stream.SendErr(err)
		}
	}()

	return stream, nil
}
