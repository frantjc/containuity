package step

import (
	api "github.com/frantjc/sequence/api/v1/step"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/services"
	"google.golang.org/grpc"
)

func NewService(opts ...StepOpt) (StepService, error) {
	svc := &stepServer{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

type stepServer struct {
	api.UnimplementedStepServer
	client *stepClient
}

type StepService interface {
	api.StepServer
	services.Service
}

var _ StepService = &stepServer{}

func (s *stepServer) RunStep(in *api.RunStepRequest, stream api.Step_RunStepServer) error {
	clientStream, err := s.client.RunStep(stream.Context(), in)
	if err != nil {
		return err
	}

	stdout, stderr := grpcio.NewLogStreamMultiplexWriter(stream)
	return grpcio.DemultiplexLogStream(clientStream, stdout, stderr)
}

func (s *stepServer) Client() (interface{}, error) {
	return s.client, nil
}

func (s *stepServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterStepServer(r, s)
}
