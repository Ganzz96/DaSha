package agent

import (
	"net"

	"github.com/ganzz96/dasha-agent/internal/log"
)

type DataMonitor struct {
	logger *log.Logger
	socket *net.UDPConn
}

func NewDataMonitor(logger *log.Logger, socket *net.UDPConn) *DataMonitor {
	return &DataMonitor{logger: logger, socket: socket}
}

func (dm *DataMonitor) Up() {
	buffer := make([]byte, 1024)

	for {
		nn, _, err := dm.socket.ReadFromUDP(buffer)
		if err != nil {
			dm.logger.Sugar().Error("Failed to read data chunk: ", err)
			continue
		}

		chunk := buffer[:nn]
		dm.logger.Sugar().Info("Get data chunk: ", string(chunk))
	}
}
