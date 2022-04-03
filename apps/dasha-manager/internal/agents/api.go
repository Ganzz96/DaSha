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
	var agent Agent

	if err := fromBody(r.Body, &agent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ac.Register(&agent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fromBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return errors.WithStack(decoder.Decode(&dest))
}
