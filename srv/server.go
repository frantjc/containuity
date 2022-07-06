package srv

import (
	"context"
	"net/http"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/sqnc"
	"github.com/frantjc/sequence/svc"
)

type Server struct {
	mux     *http.ServeMux
	runtime runtime.Runtime
}

func NewHandler(ctx context.Context, opts ...Opt) (*Server, error) {
	var (
		server = &Server{
			mux: http.NewServeMux(),
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
		func() (string, http.Handler) {
			return sqnc.NewRuntimeServiceHandler(&svc.SqncRuntimeServiceHandler{
				Runtime: server.runtime,
			})
		},
	} {
		path, handler := f()
		server.mux.Handle(path, handler)
	}

	return server, nil
}

func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
