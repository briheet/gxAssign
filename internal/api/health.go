package api

import (
	"net/http"

	"go.uber.org/zap"
)

func (a *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collections := a.database.Collection("health")

	_, err := collections.InsertOne(ctx, map[string]interface{}{"status": "healthy"})
	if err != nil {
		a.logger.Error("Failed to insert health check document", zap.Error(err))
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	a.logger.Info("healthCheckHandler write successful", zap.Any("health", "ok"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Service is healthy"))
}
