package filestream

import (
	"github.com/ganzz96/dasha-manager/internal/agents"
	"github.com/pkg/errors"
)

type agentResolver interface {
	GetAgent(id string) (*agents.Agent, error)
}

type FilestreamController struct {
	agentResolver agentResolver
}

func New(agentResolver agentResolver) *FilestreamController {
	return &FilestreamController{agentResolver: agentResolver}
}

func (fc *FilestreamController) Upload(info *UploadRequest) (UploadResponse, error) {
	agent, err := fc.agentResolver.GetAgent(info.AgentID)
	if err != nil {
		return UploadResponse{}, errors.WithStack(err)
	}

	return UploadResponse{Conn: agent.Conn}, nil
}
