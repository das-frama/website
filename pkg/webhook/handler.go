package webhook

import (
	"fmt"
	"net/http"
)

// Handler provides http handlers for post model.
type Handler interface {
	Update(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	gitService Service
}

type payload struct {
}

// NewHandler creates a new handler for webhook.
func NewHandler(gitService Service) Handler {
	return &handler{gitService}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	if err := h.gitService.Proccess(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Alles gut.")
}
