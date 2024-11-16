package api

import (
	"encoding/json"
	"net/http"

	"github.com/briheet/gxAssign/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (a *api) login(w http.ResponseWriter, r *http.Request) {
}

func (a *api) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.logger.Error("error via decoding User details in register", zap.Error(err))
		http.Error(w, "invaild request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if user.Email == "" || user.Password == "" {
		a.logger.Error("userEmail or userPassword issue in register", zap.Error(err))
		http.Error(w, "Email or Password not found", http.StatusBadRequest)
		return
	}

	filter := bson.M{"email": user.Email}
	var existingUser models.User

	err = a.database.Collection("users").FindOne(ctx, filter).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Error("error checking for user in register", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if existingUser.Email != "" {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	_, err = a.database.Collection("users").InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		a.logger.Error("Error creating a user", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
}
