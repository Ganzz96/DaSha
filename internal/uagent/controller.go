package uagent

import (
	"fmt"
	"net"

	"github.com/ccding/go-stun/stun"
	"github.com/ganzz96/dasha/internal/uagent/clients"
	"github.com/pkg/errors"
)

var (
	ErrSymmetricNATUnsupported = errors.New("Symmetric NAT unsupported")
)

type agentManager interface {
	ExchangeExternalAddr(req clients.ExchangeRequest) (clients.ExchangeResponse, error)
}

type agentGateway interface {
	PingNAgent(nAgentAddr *net.UDPAddr) error
}

type Controller struct {
	am           agentManager
	agw          agentGateway
	externalAddr *stun.Host
}

func New(natType stun.NATType, externalAddr *stun.Host, am agentManager, agw agentGateway) (*Controller, error) {
	if natType == stun.NATSymmetric {
		return nil, errors.WithStack(ErrSymmetricNATUnsupported)
	}

	return &Controller{am: am, agw: agw, externalAddr: externalAddr}, nil
}

// nagentID is temporal debug field
func (c *Controller) UploadFile(name string, path string, nagentID string) error {
	fmt.Println("Upload file", name, path, nagentID)

	agentDescription, err := c.am.ExchangeExternalAddr(clients.ExchangeRequest{
		UAgentAddr: c.externalAddr.String(),
		NAgentID:   nagentID,
		UAgentID:   "user_agent_id",
	})
	if err != nil {
		fmt.Println("Failed to get agent description", err)
		return errors.WithStack(err)
	}

	fmt.Println("Got agent description", agentDescription)

	agentAddr, err := net.ResolveUDPAddr("udp", agentDescription.NAgentExternalAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.agw.PingNAgent(agentAddr)
}
