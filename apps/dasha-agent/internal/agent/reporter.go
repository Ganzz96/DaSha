package agent

import (
	"time"

	"github.com/ganzz96/dasha-agent/internal/config"
	"github.com/ganzz96/dasha-agent/internal/log"
)

type managerClient interface {
	ReportOnline(agentID string) error
}

type AgentReporter struct {
	logger         *log.Logger
	managerClient  managerClient
	reportInterval time.Duration
	meta           *config.AgentMeta
}

func NewAgentReporter(logger *log.Logger, managerClient managerClient, meta *config.AgentMeta, reportInterval time.Duration) *AgentReporter {
	return &AgentReporter{
		logger:         logger,
		managerClient:  managerClient,
		reportInterval: reportInterval,
		meta:           meta,
	}
}

func (ar *AgentReporter) Up() {
	ticker := time.NewTicker(ar.reportInterval)

	for {
		if err := ar.managerClient.ReportOnline(ar.meta.AgentID); err != nil {
			ar.logger.Sugar().Error(err)
		}

		<-ticker.C
	}
}
