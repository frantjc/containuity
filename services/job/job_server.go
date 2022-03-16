package job

import (
	api "github.com/frantjc/sequence/api/v1/job"
	"github.com/frantjc/sequence/internal/grpcio"
	"github.com/frantjc/sequence/services"
	"google.golang.org/grpc"
)

func NewService(opts ...JobOpt) (JobService, error) {
	svc := &jobServer{}
	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

type jobServer struct {
	api.UnimplementedJobServer
	client *jobClient
}

type JobService interface {
	api.JobServer
	services.Service
}

var _ JobService = &jobServer{}

func (s *jobServer) RunJob(in *api.RunJobRequest, stream api.Job_RunJobServer) error {
	clientStream, err := s.client.RunJob(stream.Context(), in)
	if err != nil {
		return err
	}

	var (
		stdout = grpcio.NewLogStreamWriter(stream)
		stderr = stdout
	)
	return grpcio.MultiplexLogStream(clientStream, stdout, stderr)
}

func (s *jobServer) Client() (interface{}, error) {
	return s.client, nil
}

func (s *jobServer) Register(r grpc.ServiceRegistrar) {
	api.RegisterJobServer(r, s)
}
