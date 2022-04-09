package agent

import (
	"encoding/json"
	"fmt"
	"net"
)

type PushAgentMessage struct {
	UAgentAddr string `json:"uagent_addr"`
	UAgentID   string `json:"uagent_id"`
}

type agentGateway interface {
	PingUAgent(uAgentAddr *net.UDPAddr) error
}

type ConnManager struct {
	socket       *net.TCPConn
	agentGateway agentGateway
}

func NewConnManager(socket *net.TCPConn, agentGateway agentGateway) *ConnManager {
	return &ConnManager{socket: socket, agentGateway: agentGateway}
}

func (cm *ConnManager) Serve() {
	buffer := make([]byte, 1024)

	for {
		nn, err := cm.socket.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read data chunk: ", err)
			continue
		}

		var push PushAgentMessage
		if err := json.Unmarshal(buffer[:nn], &push); err != nil {
			fmt.Println("Failed to unmarshall data: ", err, " push: ", string(buffer[:nn]))
			continue
		}

		uAgentAddr, err := net.ResolveUDPAddr("udp", push.UAgentAddr)
		if err != nil {
			fmt.Println("Failed to resolve tcp addr: ", err, " addr: ", push.UAgentAddr)
			continue
		}

		if err := cm.agentGateway.PingUAgent(uAgentAddr); err != nil {
			fmt.Println("Failed to process agent push: ", err, " push: ", string(buffer[:nn]))
			continue
		}

		fmt.Println("Get push: ", push)
	}
}
