package clients

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	uploadEndpoint = "/upload"
)

type UploadRequest struct {
	AgentID string `json:"agent_id"`
}

type UploadResponse struct {
	Conn string `json:"conn"`
}

type DashaManagerClient struct {
	http     *http.Client
	endpoint string
}

func NewDashaManagerClient(endpoint string) (*DashaManagerClient, error) {
	return &DashaManagerClient{
		http:     &http.Client{},
		endpoint: endpoint,
	}, nil
}

func (c *DashaManagerClient) Upload(upload UploadRequest) (UploadResponse, error) {
	reqBody, err := json.Marshal(upload)
	if err != nil {
		return UploadResponse{}, errors.WithStack(err)
	}

	url, err := buildURL(c.endpoint, uploadEndpoint)
	if err != nil {
		return UploadResponse{}, errors.WithStack(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return UploadResponse{}, errors.WithStack(err)
	}

	rawResp, err := c.http.Do(req)
	if err != nil {
		return UploadResponse{}, errors.WithStack(err)
	}
	if is4xx(rawResp.StatusCode) || is5xx(rawResp.StatusCode) {
		return UploadResponse{}, errors.WithMessagef(BadResponseStatusCode, "Status Code: %v", rawResp.StatusCode)
	}

	var resp UploadResponse

	err = fromBody(rawResp.Body, &resp)
	if err != nil {
		return UploadResponse{}, errors.WithStack(err)
	}

	return resp, nil
}
