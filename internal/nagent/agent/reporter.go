package agent

import (
	"fmt"
	"time"

	"github.com/ganzz96/dasha/internal/nagent/clients"
	"github.com/ganzz96/dasha/internal/nagent/config"
)

type managerClient interface {
	ReportOnline(req *clients.ReportRequest) error
}

type AgentReporter struct {
	managerClient  managerClient
	reportInterval time.Duration
	meta           *config.AgentMeta
	nAgentAddr     string
}

func NewAgentReporter(managerClient managerClient, meta *config.AgentMeta, nAgentAddr string, reportInterval time.Duration) *AgentReporter {
	return &AgentReporter{
		managerClient:  managerClient,
		reportInterval: reportInterval,
		nAgentAddr:     nAgentAddr,
		meta:           meta,
	}
}

func (ar *AgentReporter) Up() {
	ticker := time.NewTicker(ar.reportInterval)

	for {
		if err := ar.managerClient.ReportOnline(&clients.ReportRequest{
			NAgentID:   ar.meta.AgentID,
			NAgentAddr: ar.nAgentAddr,
		}); err != nil {
			fmt.Println(err)
		}

		<-ticker.C
	}
}
