package clients

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	maxPingAttempts = 10
)

var (
	ErrPingAgentUnavailable = errors.New("Agent does not send pong")
)

type AgentGateway struct {
	socket      *net.UDPConn
	pongStreams sync.Map
}

func NewAgentGateway(socket *net.UDPConn) *AgentGateway {
	return &AgentGateway{
		socket:      socket,
		pongStreams: sync.Map{},
	}
}

func (agw *AgentGateway) Serve() {
	for {
		buffer := make([]byte, 1024)
		nn, from, err := agw.socket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Failed to read udp msg from", err)
			continue
		}
		fmt.Println("Got msg on udp port", string(buffer[:nn]), " from ", from.String())

		if load, ok := agw.pongStreams.Load(from.String()); ok {
			stream := load.(chan []byte)
			stream <- buffer[:nn]
		}
	}
}

func (agw *AgentGateway) PingNAgent(nAgentAddr *net.UDPAddr) error {
	stream := agw.openOrCreatePongStream(nAgentAddr)
	ticker := time.NewTicker(time.Second)
	attempts := 0

	for {
		select {
		case <-ticker.C:
			if attempts > maxPingAttempts {
				return errors.WithStack(ErrPingAgentUnavailable)
			}

			attempts++

			_, err := agw.socket.WriteToUDP([]byte("Hello from UAGENT"), nAgentAddr)
			if err != nil {
				return errors.WithStack(err)
			}

			fmt.Println("Ping msg: ", "Hello from UAGENT", " to udp address ", nAgentAddr.String())
		case msg := <-stream:
			fmt.Println("Pong  msg: ", msg, " from udp address ", nAgentAddr.String())
			return nil
		}
	}
}

func (agw *AgentGateway) openOrCreatePongStream(nAgentAddr *net.UDPAddr) chan []byte {
	stream := make(chan []byte)

	load, loaded := agw.pongStreams.LoadOrStore(nAgentAddr.String(), stream)
	if loaded {
		stream = load.(chan []byte)
	}

	return stream
}
