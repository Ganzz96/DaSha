package agent_manager

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type storage interface {
	GetAgent(id string) (*Agent, error)
	CreateAgent(agent *Agent) error
	UpdateAgentConn(id string, conn string) error
}

type agentGateway interface {
	PushAgent(nAgentAddr string, msg PushAgentMessage) error
}

type AgentController struct {
	db       storage
	gwPusher agentGateway
}

func New(storage storage, gwPusher agentGateway) *AgentController {
	return &AgentController{db: storage, gwPusher: gwPusher}
}

func (am *AgentController) Exchange(req *PostExchangeRequest) (PostExchangeResponse, error) {
	fmt.Println("Exchange uagent addr", req)

	agent, err := am.db.GetAgent(req.NAgentID)
	if err != nil {
		return PostExchangeResponse{}, errors.WithStack(err)
	}

	fmt.Println("Load agent", agent)

	if agent == nil {
		return PostExchangeResponse{}, errors.WithStack(ErrAgentDoesNotExist)
	}
	if agent.Conn == nil {
		return PostExchangeResponse{}, errors.WithStack(ErrAgentConnEmpty)
	}

	if err := am.gwPusher.PushAgent(*agent.Conn, PushAgentMessage{
		UAgentAddr: req.UAgentAddr,
		UAgentID:   req.UAgentID,
	}); err != nil {
		return PostExchangeResponse{}, errors.WithStack(err)
	}

	return PostExchangeResponse{NAgentExternalAddr: *agent.Conn}, nil

}

func (am *AgentController) Register(req *RegisterAgentRequest) (RegisterAgentResponse, error) {
	now := time.Now()

	agent := Agent{
		ID:        uuid.New().String(),
		CreatedAt: &now,
	}

	if err := am.db.CreateAgent(&agent); err != nil {
		return RegisterAgentResponse{}, errors.WithStack(err)
	}

	return RegisterAgentResponse{
		AgentID: agent.ID,
	}, nil
}

func (am *AgentController) Report(req *PostReportRequest) error {
	existed, err := am.db.GetAgent(req.NAgentID)
	if err != nil {
		return errors.WithStack(err)
	}

	if existed == nil {
		return errors.WithStack(ErrAgentDoesNotExist)
	}

	return am.db.UpdateAgentConn(req.NAgentID, req.NAgentAddr)
}
