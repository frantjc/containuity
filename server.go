package sequence

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/services"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer  *grpc.Server
	httpHandler http.Handler
}

var _ http.Handler = &Server{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(
		r.Header.Get("Content-Type"), "application/grpc") {
		s.grpcServer.ServeHTTP(w, r)
	} else {
		s.httpHandler.ServeHTTP(w, r)
	}
}

func (s *Server) Serve(l net.Listener) error {
	return s.grpcServer.Serve(l)
}

func NewServer(ctx context.Context, opts ...ServerOpt) (*Server, error) {
	so := &serverOpts{}
	for _, opt := range opts {
		err := opt(so)
		if err != nil {
			return nil, err
		}
	}

	var (
		grpcServer = grpc.NewServer()
		runtime    = so.runtime
		svcOpts    = []services.Opt{services.WithRuntime(runtime)}
	)
	imageService, err := services.NewImageService(svcOpts...)
	if err != nil {
		return nil, err
	}
	imageService.Register(grpcServer)
	log.Info("registered image service")

	containerService, err := services.NewContainerService(svcOpts...)
	if err != nil {
		return nil, err
	}
	containerService.Register(grpcServer)
	log.Info("registered container service")

	volumeService, err := services.NewVolumeService(svcOpts...)
	if err != nil {
		return nil, err
	}
	volumeService.Register(grpcServer)
	log.Info("registered volume service")

	stepService, err := services.NewStepService(svcOpts...)
	if err != nil {
		return nil, err
	}
	stepService.Register(grpcServer)
	log.Info("registered step service")

	jobService, err := services.NewJobService(svcOpts...)
	if err != nil {
		return nil, err
	}
	jobService.Register(grpcServer)
	log.Info("registered job service")

	workflowService, err := services.NewWorkflowService(svcOpts...)
	if err != nil {
		return nil, err
	}
	workflowService.Register(grpcServer)
	log.Info("registered workflow service")

	return &Server{
		grpcServer: grpcServer,
	}, nil
}
