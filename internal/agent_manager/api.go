package agent_manager

import (
	"errors"
	"net/http"

	"github.com/ganzz96/dasha/internal/common/utils"
	"github.com/go-chi/chi"
)

type RegisterAgentRequest struct {
}

type RegisterAgentResponse struct {
	AgentID string `json:"agent_id"`
}

type PostExchangeRequest struct {
	UAgentAddr string `json:"uagent_addr"`
	NAgentID   string `json:"nagent_id"`
	UAgentID   string `json:"uagent_id"`
}

type PostExchangeResponse struct {
	NAgentExternalAddr string `json:"nagent_external_addr"`
}

type PostReportRequest struct {
	NAgentAddr string `json:"nagent_addr"`
	NAgentID   string `json:"nagent_id"`
}

func (ac *AgentController) RegisterAPI(router *chi.Mux) {
	router.Post("/agents", ac.registerAgent)
	router.Post("/exchange", ac.postExchange)
	router.Post("/report", ac.postReport)
}

func (ac *AgentController) postExchange(w http.ResponseWriter, r *http.Request) {
	var req PostExchangeRequest

	if err := utils.FromBody(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := ac.Exchange(&req)
	if err != nil {
		if errors.Is(err, ErrAgentDoesNotExist) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := utils.ToBody(w, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ac *AgentController) registerAgent(w http.ResponseWriter, r *http.Request) {
	var req RegisterAgentRequest

	if err := utils.FromBody(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := ac.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := utils.ToBody(w, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ac *AgentController) postReport(w http.ResponseWriter, r *http.Request) {
	var req PostReportRequest

	if err := utils.FromBody(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ac.Report(&req); err != nil {
		if errors.Is(err, ErrAgentDoesNotExist) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
