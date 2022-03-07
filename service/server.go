package service

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/frantjc/sequence/service/workflow"
	"google.golang.org/grpc"
)

type Server interface {
	http.Handler
	Serve(l net.Listener) error
}

type server struct {
	grpcsrv  *grpc.Server
	httphdlr http.Handler
}

var _ Server = &server{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(
		r.Header.Get("Content-Type"), "application/grpc") {
		s.grpcsrv.ServeHTTP(w, r)
	} else {
		s.httphdlr.ServeHTTP(w, r)
	}
}

func (s *server) Serve(l net.Listener) error {
	return s.grpcsrv.Serve(l)
}

func New(ctx context.Context, opts ...ServerOpt) (Server, error) {
	so := &serverOpts{}
	for _, opt := range opts {
		err := opt(so)
		if err != nil {
			return nil, err
		}
	}

	s := grpc.NewServer()

	ropts := []workflow.WorkflowOpt{}
	if so.conf != nil {
		ropts = append(ropts, workflow.WithConfig(so.conf))
	}
	workflow.Register(ctx, s, ropts...)

	return &server{
		grpcsrv: s,
	}, nil
}
