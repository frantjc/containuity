package workflow

import (
	api "github.com/frantjc/sequence/api/v1/step"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/workflow"
)

func (s *workflowServer) RunStep(req *api.RunStepRequest, stream api.Step_RunStepServer) error {
	opts := []workflow.RunOpt{
		workflow.WithStdout(grpcio.NewLogStreamWriter(stream)),
		workflow.WithGitHubToken(s.conf.GitHub.Token),
	}

	if s.conf.Verbose {
		opts = append(opts, workflow.WithVerbose)
	}

	err := workflow.RunStep(stream.Context(), s.r, convert.TypeStepToStep(req.Step), opts...)
	if err != nil {
		return err
	}

	return nil
}
