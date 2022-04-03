package monitor

import (
	"net"
	"time"

	"github.com/ganzz96/dasha-manager/internal/log"
)

type agentReporter interface {
	Report(aID string, conn string) error
}

type AgentMonitor struct {
	logger        *log.Logger
	agentReporter agentReporter
}

func New(logger *log.Logger, agentReporter agentReporter) *AgentMonitor {
	return &AgentMonitor{logger: logger, agentReporter: agentReporter}
}

func (am *AgentMonitor) Serve(host string, port int) {
	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	})
	if err != nil {
		am.logger.Sugar().Fatal(err)
	}
	defer socket.Close()

	for {
		time.Sleep(time.Second)

		buffer := make([]byte, 100)
		nn, agentUDPAddr, err := socket.ReadFromUDP(buffer)
		if err != nil {
			am.logger.Sugar().Error(err)
			continue
		}

		if nn == 0 && len(buffer) == 0 {
			continue
		}

		if err = am.agentReporter.Report(string(buffer[:nn]), agentUDPAddr.String()); err != nil {
			am.logger.Sugar().Error(err)
			continue
		}
	}
}
