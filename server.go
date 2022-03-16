package sequence

import (
	"context"
	"net"
	"net/http"
	"strings"

	containerapi "github.com/frantjc/sequence/api/v1/container"
	imageapi "github.com/frantjc/sequence/api/v1/image"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/services/container"
	"github.com/frantjc/sequence/services/image"
	"github.com/frantjc/sequence/services/job"
	"github.com/frantjc/sequence/services/step"
	"github.com/frantjc/sequence/services/workflow"
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
	)
	imageService, err := image.NewService(image.WithRuntime(runtime))
	if err != nil {
		return nil, err
	}
	imageService.Register(grpcServer)
	log.Info("registered image service")

	containerService, err := container.NewService(container.WithRuntime(runtime))
	if err != nil {
		return nil, err
	}
	containerService.Register(grpcServer)
	log.Info("registered container service")

	var (
		imageClient, _     = imageService.Client()
		containerClient, _ = containerService.Client()
	)
	if ic, ok := imageClient.(imageapi.ImageClient); ok {
		if cc, ok := containerClient.(containerapi.ContainerClient); ok {
			log.Info("using container and image services as runtime")
			runtime = NewGRPCRuntime(ic, cc)
		}
	}

	stepService, err := step.NewService(step.WithRuntime(runtime))
	if err != nil {
		return nil, err
	}
	stepService.Register(grpcServer)
	log.Info("registered step service")

	jobService, err := job.NewService(job.WithRuntime(runtime))
	if err != nil {
		return nil, err
	}
	jobService.Register(grpcServer)
	log.Info("registered job service")

	workflowService, err := workflow.NewService(workflow.WithRuntime(runtime))
	if err != nil {
		return nil, err
	}
	workflowService.Register(grpcServer)
	log.Info("registered workflow service")

	return &Server{
		grpcServer: grpcServer,
	}, nil
}
