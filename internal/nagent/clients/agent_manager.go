package clients

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ganzz96/dasha/internal/common/utils"
	"github.com/pkg/errors"
)

const (
	agentEndpoint  = "/agents"
	reportEndpoint = "/report"
)

var (
	BadResponseStatusCode = errors.New("bad response status code")
)

type RegisterRequest struct {
}

type RegisterResponse struct {
	AgentID string `json:"agent_id"`
}

type ReportRequest struct {
	NAgentID   string `json:"nagent_id"`
	NAgentAddr string `json:"nagent_addr"`
}

type AgentManagerClient struct {
	http         *http.Client
	httpEndpoint string
}

func NewAgentManagerClient(endpoint string) (*AgentManagerClient, error) {
	return &AgentManagerClient{
		http:         &http.Client{},
		httpEndpoint: endpoint,
	}, nil
}

func (agc *AgentManagerClient) ReportOnline(req *ReportRequest) error {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return errors.WithStack(err)
	}

	url, err := utils.BuildURL(agc.httpEndpoint, reportEndpoint)
	if err != nil {
		return errors.WithStack(err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.WithStack(err)
	}

	rawResp, err := agc.http.Do(httpReq)
	if err != nil {
		return errors.WithStack(err)
	}

	if utils.Is4xx(rawResp.StatusCode) || utils.Is5xx(rawResp.StatusCode) {
		return errors.WithMessagef(BadResponseStatusCode, "Status Code: %v", rawResp.StatusCode)
	}

	return nil
}

func (agc *AgentManagerClient) RegisterAgent(req *RegisterRequest) (RegisterResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return RegisterResponse{}, errors.WithStack(err)
	}

	url, err := utils.BuildURL(agc.httpEndpoint, agentEndpoint)
	if err != nil {
		return RegisterResponse{}, errors.WithStack(err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return RegisterResponse{}, errors.WithStack(err)
	}

	rawResp, err := agc.http.Do(httpReq)
	if err != nil {
		return RegisterResponse{}, errors.WithStack(err)
	}

	if utils.Is4xx(rawResp.StatusCode) || utils.Is5xx(rawResp.StatusCode) {
		return RegisterResponse{}, errors.WithMessagef(BadResponseStatusCode, "Status Code: %v", rawResp.StatusCode)
	}

	var resp RegisterResponse

	err = utils.FromBody(rawResp.Body, &resp)
	if err != nil {
		return RegisterResponse{}, errors.WithStack(err)
	}

	return resp, nil
}
