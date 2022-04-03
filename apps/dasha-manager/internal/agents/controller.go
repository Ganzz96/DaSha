package agents

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type storage interface {
	GetAgent(id string) (*Agent, error)
	CreateAgent(agent *Agent) error
	UpdateAgentConn(id string, conn string) error
}

type AgentController struct {
	db storage
}

func New(storage storage) *AgentController {
	return &AgentController{db: storage}
}

func (am *AgentController) Register(req *RegisterRequest) (RegisterResponse, error) {
	now := time.Now()

	agent := Agent{
		ID:        uuid.New().String(),
		CreatedAt: &now,
	}

	if err := am.db.CreateAgent(&agent); err != nil {
		return RegisterResponse{}, errors.WithStack(err)
	}

	return RegisterResponse{
		AgentID: agent.ID,
	}, nil
}

func (am *AgentController) Report(aID string, conn string) error {
	existed, err := am.db.GetAgent(aID)
	if err != nil {
		return errors.WithStack(err)
	}

	if existed == nil {
		return errors.WithStack(ErrAgentDoesNotExist)
	}

	return am.db.UpdateAgentConn(aID, conn)
}

func (am *AgentController) GetAgent(id string) (*Agent, error) {
	return am.db.GetAgent(id)
}
