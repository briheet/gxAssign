package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type api struct {
	logger *zap.Logger
	client *http.Client
}

func NewAPI(context context.Context, logger *zap.Logger) *api {
	client := http.Client{}
	return &api{
		logger: logger,
		client: &client,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/health", a.healthCheckHandler).Methods("GET")
	return r
}
