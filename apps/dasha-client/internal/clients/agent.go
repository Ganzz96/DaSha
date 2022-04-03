package clients

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type AgentClient struct {
	socket *net.UDPConn
}

func NewAgentClient(agentConn string) (*AgentClient, error) {
	udpDestAddr, err := net.ResolveUDPAddr("udp4", agentConn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	socket, err := net.DialUDP("udp", nil, udpDestAddr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &AgentClient{
		socket: socket,
	}, nil
}

func (ac *AgentClient) Transmit(message string) error {
	nn, err := ac.socket.Write([]byte(message))
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("Successfully sent ", nn, " bytes to udp address ", ac.socket.RemoteAddr().String())
	return nil
}
