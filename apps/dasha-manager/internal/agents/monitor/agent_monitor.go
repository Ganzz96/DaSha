package monitor

import (
	"net"

	"github.com/ganzz96/dasha-manager/internal/log"
)

type agentReporter interface {
	Report(aID string, conn string) error
}

type AgentMonitor struct {
	logger        *log.Logger
	agentReporter agentReporter
	socket        *net.UDPConn
}

func New(logger *log.Logger, socket *net.UDPConn, agentReporter agentReporter) *AgentMonitor {
	return &AgentMonitor{logger: logger, socket: socket, agentReporter: agentReporter}
}

func (am *AgentMonitor) Serve() {
	for {
		buffer := make([]byte, 100)
		nn, agentUDPAddr, err := am.socket.ReadFromUDP(buffer)
		if err != nil {
			am.logger.Sugar().Error(err)
			continue
		}

		if nn == 0 && len(buffer) == 0 {
			continue
		}

		agentID := string(buffer[:nn])

		if err = am.agentReporter.Report(agentID, agentUDPAddr.String()); err != nil {
			am.logger.Sugar().Errorf("Failed to get report by agent with id %s. Got error %v.", agentID, err)
			continue
		}
	}
}
