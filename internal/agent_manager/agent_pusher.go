package agent_manager

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"
)

type AgentPusher struct {
	socket *net.TCPListener
	conns  sync.Map
}

type PushAgentMessage struct {
	UAgentAddr string `json:"uagent_addr"`
	UAgentID   string `json:"uagent_id"`
}

func NewAgentPusher(socket *net.TCPListener) *AgentPusher {
	return &AgentPusher{socket: socket, conns: sync.Map{}}
}

func (ap *AgentPusher) Serve() {
	for {
		conn, err := ap.socket.AcceptTCP()
		if err != nil {
			fmt.Println("Failed to accept tcp conn", err)
			continue
		}

		if err := conn.SetKeepAlive(true); err != nil {
			fmt.Println("Failed to set keep alive tcp conn", err)
			continue
		}

		nAgentAddr := conn.RemoteAddr().String()

		load, loaded := ap.conns.LoadOrStore(nAgentAddr, conn)
		if !loaded {
			continue
		}

		prevConn := load.(*net.TCPConn)
		if err := prevConn.Close(); err != nil {
			fmt.Println("Failed to close prev tcp conn", err)
			continue
		}

		ap.conns.Store(nAgentAddr, conn)
	}
}

func (ap *AgentPusher) PushAgent(nAgentAddr string, msg PushAgentMessage) error {
	load, ok := ap.conns.Load(nAgentAddr)
	if !ok {
		return errors.WithStack(ErrAgentIsNotConnected)
	}

	encoded, err := json.Marshal(msg)
	if err != nil {
		return errors.WithStack(err)
	}

	conn := load.(*net.TCPConn)
	if _, err := conn.Write(encoded); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
