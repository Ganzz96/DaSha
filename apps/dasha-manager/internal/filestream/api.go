package filestream

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func (fc *FilestreamController) RegisterAPI(router *chi.Mux) {
	router.Post("/upload", fc.upload)
}

func (fc *FilestreamController) upload(w http.ResponseWriter, r *http.Request) {
	var info UploadRequest

	if err := fromBody(r.Body, &info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := fc.Upload(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fromBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return errors.WithStack(decoder.Decode(&dest))
}
