package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ganzz96/dasha-manager/internal/gateway"
	"github.com/ganzz96/dasha-manager/internal/log"
	"github.com/ganzz96/dasha-manager/internal/storage"
)

func main() {
	router := chi.NewRouter()
	logger := log.New()

	if _, err := storage.New(logger, "sqlite.db"); err != nil {
		logger.Sugar().Panic(err)
	}

	gw := gateway.New()
	gw.RegisterAPI(router)

	if err := http.ListenAndServe(":7777", router); err != nil {
		logger.Sugar().Panic(err)
	}
}
