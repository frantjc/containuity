package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler interface {
	http.Handler
	Method() string
	Path() string
}

func NewServer() http.Handler {
	router := httprouter.New()

	for _, handler := range []Handler{
		&WebhookHandler{},
	} {
		router.Handler(
			handler.Method(),
			handler.Path(),
			handler,
		)
	}

	return router
}
