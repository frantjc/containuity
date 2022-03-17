package workflow

import (
	api "github.com/frantjc/sequence/api/v1/workflow"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/services"
	"google.golang.org/grpc"
)

func NewService(opts ...WorkflowOpt) (WorkflowService, error) {
	svc := &workflowServer{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

type workflowServer struct {
	api.UnimplementedWorkflowServer
	client *workflowClient
}

type WorkflowService interface {
	api.WorkflowServer
	services.Service
}

var _ WorkflowService = &workflowServer{}

func (s *workflowServer) RunWorkflow(in *api.RunWorkflowRequest, stream api.Workflow_RunWorkflowServer) error {
	clientStream, err := s.client.RunWorkflow(stream.Context(), in)
	if err != nil {
		return err
	}

	stdout, stderr := grpcio.NewLogStreamMultiplexWriter(stream)
	return grpcio.DemultiplexLogStream(clientStream, stdout, stderr)
}

func (s *workflowServer) Client() (interface{}, error) {
	return s.client, nil
}

func (s *workflowServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterWorkflowServer(r, s)
}
