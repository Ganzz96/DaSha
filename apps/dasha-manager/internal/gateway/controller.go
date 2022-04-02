package gateway

import (
	"net/http"

	"github.com/go-chi/chi"
)

type GatewayController struct {
}

func New() *GatewayController {
	return &GatewayController{}
}

func (gw *GatewayController) RegisterAPI(router *chi.Mux) {
	router.Get("/agents", gw.listAgents)
}

func (gw *GatewayController) listAgents(w http.ResponseWriter, r *http.Request) {

}
