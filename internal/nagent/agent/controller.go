package agent

import (
	"github.com/ganzz96/dasha/internal/nagent/clients"
	"github.com/pkg/errors"
)

type registrator interface {
	RegisterAgent(req *clients.RegisterRequest) (clients.RegisterResponse, error)
}

type agentMetaSaver interface {
	SyncAgentID(agID string) error
}

type AgentController struct {
	registrator    registrator
	agentMetaSaver agentMetaSaver
}

func New(registrator registrator, agentMetaSaver agentMetaSaver) *AgentController {
	return &AgentController{registrator: registrator, agentMetaSaver: agentMetaSaver}
}

func (ac *AgentController) Register() error {
	resp, err := ac.registrator.RegisterAgent(&clients.RegisterRequest{})
	if err != nil {
		return errors.WithStack(err)
	}

	return ac.agentMetaSaver.SyncAgentID(resp.AgentID)
}
