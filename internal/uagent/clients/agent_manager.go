package clients

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ganzz96/dasha/internal/common/utils"
	"github.com/pkg/errors"
)

const (
	exchangeEndpoint = "/exchange"
)

var (
	BadResponseStatusCode = errors.New("bad response status code")
)

type ExchangeRequest struct {
	UAgentAddr string `json:"uagent_addr"`
	NAgentID   string `json:"nagent_id"`
	UAgentID   string `json:"uagent_id"`
}

type ExchangeResponse struct {
	NAgentExternalAddr string `json:"nagent_external_addr"`
}

type AgentManagerClient struct {
	http     *http.Client
	endpoint string
}

func NewAgentManagerClient(endpoint string) (*AgentManagerClient, error) {
	return &AgentManagerClient{
		http:     &http.Client{},
		endpoint: endpoint,
	}, nil
}

func (c *AgentManagerClient) ExchangeExternalAddr(req ExchangeRequest) (ExchangeResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return ExchangeResponse{}, errors.WithStack(err)
	}

	url, err := utils.BuildURL(c.endpoint, exchangeEndpoint)
	if err != nil {
		return ExchangeResponse{}, errors.WithStack(err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return ExchangeResponse{}, errors.WithStack(err)
	}

	rawResp, err := c.http.Do(httpReq)
	if err != nil {
		return ExchangeResponse{}, errors.WithStack(err)
	}
	if utils.Is4xx(rawResp.StatusCode) || utils.Is5xx(rawResp.StatusCode) {
		return ExchangeResponse{}, errors.WithMessagef(BadResponseStatusCode, "Status Code: %v", rawResp.StatusCode)
	}

	var resp ExchangeResponse

	err = utils.FromBody(rawResp.Body, &resp)
	if err != nil {
		return ExchangeResponse{}, errors.WithStack(err)
	}

	return resp, nil
}
