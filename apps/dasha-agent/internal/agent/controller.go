package agent

import "github.com/pkg/errors"

type registrator interface {
	RegisterAgent() (string, error)
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
	agID, err := ac.registrator.RegisterAgent()
	if err != nil {
		return errors.WithStack(err)
	}

	return ac.agentMetaSaver.SyncAgentID(agID)
}
