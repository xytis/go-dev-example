package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/xytis/go-dev-example/internal/queue"
	"github.com/xytis/go-dev-example/internal/x"
)

type handler struct {
	queue queue.SimpleQueue
}

func NewHandler(queue queue.SimpleQueue) http.Handler {
	h := handler{
		queue: queue,
	}

	mux := http.NewServeMux()
	mux.Handle("POST /", http.HandlerFunc(h.IngestMessage))
	mux.Handle("GET /", http.HandlerFunc(h.RetrieveMessage))

	return mux
}

func (h *handler) IngestMessage(w http.ResponseWriter, r *http.Request) {
	// Extract message from HTTP
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		x.ApiError(w, err)
		return
	}

	// Sanitize and parse
	parsed, err := queue.ParseMessage(string(msg))
	if err != nil {
		x.ApiError(w, err)
		return
	}

	// Do the business operation
	err = h.queue.Push(parsed)
	if err != nil {
		x.ApiError(w, err)
	}

	// Produce protocol response
	x.ApiSuccess(w)
}

func (h *handler) RetrieveMessage(w http.ResponseWriter, _ *http.Request) {
	msg, err := h.queue.Pull()
	if errors.Is(err, queue.EOQueue) {
		x.ApiEmpty(w)
		return
	}
	if err != nil {
		x.ApiError(w, err)
		return
	}

	// Convert into protocol
	payload := []byte(msg)

	x.ApiPayload(w, payload)
}
