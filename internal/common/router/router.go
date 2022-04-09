package router

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/ganzz96/dasha/internal/common/log"
)

func NewRouter(logger *log.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.RequestID)

	//router.Use(middleware.RequestLogger(...))

	return router
}
