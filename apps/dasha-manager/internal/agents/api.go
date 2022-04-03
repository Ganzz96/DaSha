package agents

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func (ac *AgentController) RegisterAPI(router *chi.Mux) {
	router.Post("/agents", ac.registerAgent)
}

func (ac *AgentController) registerAgent(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := fromBody(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := ac.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := toBody(w, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fromBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return errors.WithStack(decoder.Decode(&dest))
}

func toBody(body io.Writer, source interface{}) error {
	encoder := json.NewEncoder(body)
	return errors.WithStack(encoder.Encode(source))
}
