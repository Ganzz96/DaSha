package uagent

import (
	"net/http"

	"github.com/ganzz96/dasha/internal/common/utils"
	"github.com/go-chi/chi"
)

type PostUploadRequest struct {
	FileName   string `json:"file_name"`
	SourcePath string `json:"source_path"`
	NAgentID   string `json:"nagent_id"` // debug field
}

type PostUploadResponse struct {
}

func (c *Controller) RegisterAPI(router *chi.Mux) {
	router.Post("/upload", c.postUpload)
}

func (c *Controller) postUpload(w http.ResponseWriter, r *http.Request) {
	var req PostUploadRequest

	if err := utils.FromBody(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.UploadFile(req.FileName, req.SourcePath, req.NAgentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
