package srv

import (
	"context"
	"net/http"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/sqnc"
	"github.com/frantjc/sequence/svc"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	mux     *http.ServeMux
	runtime runtime.Runtime
}

func NewHandler(ctx context.Context, opts ...Opt) (http.Handler, error) {
	var (
		server = &Server{
			mux: http.NewServeMux(),
		}
		err error
	)
	for _, opt := range opts {
		if err := opt(ctx, server); err != nil {
			return nil, err
		}
	}

	if server.runtime == nil {
		// get any runtime, starting with one
		// specified by SQNC_RUNTIME
		if server.runtime, err = runtimes.GetRuntime(ctx); err != nil {
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

	return h2c.NewHandler(server, &http2.Server{}), nil
}

func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
