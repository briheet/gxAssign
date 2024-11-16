package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/briheet/gxAssign/internal/db"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type api struct {
	logger      *zap.Logger
	client      *http.Client
	mongoClient *mongo.Client
	database    *mongo.Database
}

func NewAPI(context context.Context, logger *zap.Logger, dbName string) *api {
	db := db.NewDB(context, dbName)
	client := http.Client{}
	return &api{
		logger:      logger,
		client:      &client,
		mongoClient: db.Client,
		database:    db.Database,
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
	r.HandleFunc("/v1/register", a.register).Methods("POST")
	r.HandleFunc("/v1/login", a.login).Methods("POST")
	return r
}
