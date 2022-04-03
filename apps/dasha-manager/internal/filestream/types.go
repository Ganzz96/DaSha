package filestream

type UploadRequest struct {
	AgentID string `json:"agent_id"`
}

type UploadResponse struct {
	Conn string `json:"conn"`
}
