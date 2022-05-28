package sequence

import (
	"context"
	"net"
	"net/http"

	"github.com/frantjc/sequence/runtime"
	_ "github.com/frantjc/sequence/runtime/docker"
	"github.com/frantjc/sequence/services"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v44/github"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer  *grpc.Server
	httpHandler http.Handler
}

var _ http.Handler = &Server{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpHandler.ServeHTTP(w, r)
}

func (s *Server) ServeGRPC(l net.Listener) error {
	return s.grpcServer.Serve(l)
}

func NewServer(ctx context.Context, runtime runtime.Runtime, opts ...ServerOpt) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)

	var (
		sOpts       = &serverOpts{}
		grpcServer  = grpc.NewServer()
		httpHandler = gin.New()
	)
	for _, opt := range opts {
		err := opt(sOpts)
		if err != nil {
			return nil, err
		}
	}

	imageService, err := services.NewImageService(runtime)
	if err != nil {
		return nil, err
	}
	imageService.Register(grpcServer)

	containerService, err := services.NewContainerService(runtime)
	if err != nil {
		return nil, err
	}
	containerService.Register(grpcServer)

	volumeService, err := services.NewVolumeService(runtime)
	if err != nil {
		return nil, err
	}
	volumeService.Register(grpcServer)

	stepService, err := services.NewStepService(runtime)
	if err != nil {
		return nil, err
	}
	stepService.Register(grpcServer)

	jobService, err := services.NewJobService(runtime)
	if err != nil {
		return nil, err
	}
	jobService.Register(grpcServer)

	workflowService, err := services.NewWorkflowService(runtime)
	if err != nil {
		return nil, err
	}
	workflowService.Register(grpcServer)

	httpHandler.GET("/readyz", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	httpHandler.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	api := httpHandler.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/github", func(ctx *gin.Context) {
				payload, err := github.ValidatePayload(ctx.Request, sOpts.webhookSecretKey)
				if err != nil {
					ctx.AbortWithStatus(500)
				}

				var _ = payload
				ctx.AbortWithStatus(200)
			})
		}
	}

	return &Server{grpcServer, httpHandler}, nil
}
