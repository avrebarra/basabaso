package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/avrebarra/basabaso/pkg/ctxops"
	"github.com/rs/zerolog/log"
)

const (
	ctxkeyrespdata = "*ctxkey:server:respdata::"
)

func sendJSON(ctx context.Context, w http.ResponseWriter, data interface{}) (err error) {
	// save necessary context data
	ctxops.SetData(ctx, ctxkeyrespdata, data)

	// build data
	bts, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("severe error: handler/sendJSON")
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(bts); err != nil {
		return
	}

	return
}

// ***

type Handler struct{ Server *Server }

func (h *Handler) HealthCheck() http.HandlerFunc {
	make := h.Server.buildHandler

	type RequestData struct {
	}

	type ResponseData struct {
	}

	return make(func(w http.ResponseWriter, r *http.Request) (err error) {
		ctx := r.Context()

		// prep and validate

		// perform operations

		// build response
		out := RespPresets[RCSuccess]
		return sendJSON(ctx, w, out.Normalize())
	})
}
