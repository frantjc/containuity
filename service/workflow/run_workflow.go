package workflow

import (
	api "github.com/frantjc/sequence/api/v1/workflow"
	"github.com/frantjc/sequence/internal/convert"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/workflow"
)

func (s *workflowServer) RunWorkflow(req *api.RunWorkflowRequest, stream api.Workflow_RunWorkflowServer) error {
	opts := []workflow.RunOpt{
		workflow.WithStdout(grpcio.NewLogStreamWriter(stream)),
		workflow.WithGitHubToken(s.conf.GitHub.Token),
		workflow.WithDefaultImage(s.conf.Runtime.Image),
	}

	if s.conf.Verbose {
		opts = append(opts, workflow.WithVerbose)
	}

	err := workflow.RunWorkflow(stream.Context(), s.r, convert.TypeWorkflowToWorkflow(req.Workflow), opts...)
	if err != nil {
		return err
	}

	return nil
}
