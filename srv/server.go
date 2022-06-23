package srv

import (
	"context"
	"net/http"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/svc"
	"github.com/gorilla/mux"
)

type Server struct {
	handler http.Handler
	runtime runtime.Runtime
}

var _ http.Handler = &Server{}

func NewServer(ctx context.Context, opts ...Opt) (*Server, error) {
	var (
		router = mux.NewRouter()
		server = &Server{
			handler: router,
		}
	)
	for _, opt := range opts {
		if err := opt(ctx, server); err != nil {
			return nil, err
		}
	}

	if server.runtime == nil {
		if err := WithDefaultRuntime(ctx, server); err != nil {
			return nil, err
		}
	}

	type NewPathHandlerFunc func() (string, http.Handler)
	for _, f := range []NewPathHandlerFunc{
		func() (string, http.Handler) {
			return sequence.NewWorkflowServiceHandler(&svc.WorkflowServiceHandler{})
		},
	} {
		path, handler := f()
		router.Handle(path, handler)
	}

	return server, nil
}

func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}
