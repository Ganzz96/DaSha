package clients

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"

	"github.com/pkg/errors"

	"github.com/ganzz96/dasha-agent/internal/log"
)

const (
	agentEndpoint = "/agents"
)

type RegisterRequest struct {
}

type RegisterResponse struct {
	AgentID string `json:"agent_id"`
}

type DashaManagerClient struct {
	logger        *log.Logger
	socket        *net.UDPConn
	http          *http.Client
	serverAddress *net.UDPAddr
	httpEndpoint  string
}

func NewDashaManagerClient(logger *log.Logger, socket *net.UDPConn, managerHTTPHostport, managerUDPHostport string) (*DashaManagerClient, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", managerUDPHostport)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &DashaManagerClient{
		logger:        logger,
		socket:        socket,
		http:          &http.Client{},
		serverAddress: serverAddr,
		httpEndpoint:  managerHTTPHostport,
	}, nil
}

func (c *DashaManagerClient) ReportOnline(agentID string) error {
	nn, err := c.socket.WriteToUDP([]byte(agentID), c.serverAddress)
	if err != nil {
		return errors.WithStack(err)
	}

	c.logger.Sugar().Info("Reported agent online status. Bytes sent: ", nn)
	return nil
}

func (c *DashaManagerClient) RegisterAgent() (string, error) {
	reqBody, err := json.Marshal(RegisterRequest{})
	if err != nil {
		return "", errors.WithStack(err)
	}

	url, err := buildURL(c.httpEndpoint, agentEndpoint)
	if err != nil {
		return "", errors.WithStack(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", errors.WithStack(err)
	}

	rawResp, err := c.http.Do(req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	var resp RegisterResponse

	err = fromBody(rawResp.Body, &resp)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resp.AgentID, nil
}
