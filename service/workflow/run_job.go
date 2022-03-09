package workflow

import (
	api "github.com/frantjc/sequence/api/v1/job"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/workflow"
)

func (s *workflowServer) RunJob(req *api.RunJobRequest, stream api.Job_RunJobServer) error {
	opts := []workflow.RunOpt{
		workflow.WithStdout(grpcio.NewLogStreamWriter(stream)),
		workflow.WithGitHubToken(s.conf.GitHub.Token),
		workflow.WithDefaultImage(s.conf.Runtime.Image),
		workflow.WithWorkdir(s.conf.RootDir),
	}

	if s.conf.Verbose {
		opts = append(opts, workflow.WithVerbose)
	}

	err := workflow.RunJob(stream.Context(), s.r, convert.TypeJobToJob(req.Job), opts...)
	if err != nil {
		return err
	}

	return nil
}
