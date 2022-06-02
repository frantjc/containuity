package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/frantjc/sequence/datastore"
	"github.com/google/go-github/v45/github"
	"github.com/julienschmidt/httprouter"
)

type WebhookHandler struct {
	datastore     datastore.Datastore
	webhookSecret []byte
}

var _ http.Handler = &WebhookHandler{}

func (h *WebhookHandler) webhookFromParam() string {
	return "source"
}

func (h *WebhookHandler) Method() string {
	return http.MethodPost
}

func (h *WebhookHandler) Path() string {
	return fmt.Sprintf("/api/v1/webhook/:%s", h.webhookFromParam())
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	if !ok {
		params = []httprouter.Param{
			{
				Key:   h.webhookFromParam(),
				Value: "github",
			},
		}
	}

	h.Handle(w, r, params)
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch params.ByName(h.webhookFromParam()) {
	case "github":
		var (
			// ctx = r.Context()
			eventType = r.Header.Get("X-GitHub-Event")
			// deliveryID = r.Header.Get("X-GitHub-Delivery")
		)

		if eventType == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := github.ValidatePayload(r, h.webhookSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		event := &github.Event{}
		if err := json.Unmarshal(b, event); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
