package api

import (
	"encoding/json"
	"net/http"

	"github.com/briheet/gxAssign/internal/models"
	"github.com/briheet/gxAssign/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (a *api) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var usersDetails models.UserAuth
	err := json.NewDecoder(r.Body).Decode(&usersDetails)
	if err != nil {

		a.logger.Error("error via decoding User details in register", zap.Error(err))
		http.Error(w, "invaild request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if usersDetails.Email == "" || usersDetails.Password == "" {
		http.Error(w, "Email or Password not found", http.StatusBadRequest)
		return
	}

	filter := bson.M{"email": usersDetails.Email}
	var existingUserDetails *models.User

	err = a.database.Collection("users").FindOne(ctx, filter).Decode(&existingUserDetails)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
		} else {
			a.logger.Error("error finding user in login", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	if !utils.CheckPasswordHash(usersDetails.Password, existingUserDetails.Password) {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
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

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		a.logger.Error("error hashing password", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	_, err = a.database.Collection("users").InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		a.logger.Error("Error creating a user", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
}

func (a *api) upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userDocs models.UserDocument

	err := json.NewDecoder(r.Body).Decode(&userDocs)
	if err != nil {
		a.logger.Error("failed to decode the document deatils", zap.Error(err))
		http.Error(w, `{"error": "invaild request payload"}`, http.StatusBadRequest)
		return
	}

	if userDocs.UserID == "" || userDocs.Task == "" || userDocs.Admin == "" {
		http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
		return
	}

	filter := bson.M{"email": userDocs.UserID}
	var user bson.M
	if err = a.database.Collection("users").FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		} else {
			a.logger.Error("Database error while verifying user", zap.Error(err))
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	_, err = a.database.Collection("documents").InsertOne(ctx, userDocs)
	if err != nil {
		a.logger.Error("Failed to insert document", zap.Error(err))
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "document uploaded successfully"})
}

func (a *api) getAdmins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("userID")

	if userID == "" {
		http.Error(w, `{"error": "UserID is required"}`, http.StatusBadRequest)
		return
	}

	var userAdmins models.UserAdmins
	filter := bson.M{"userID": userID}
	err := a.database.Collection("useradmins").FindOne(ctx, filter).Decode(&userAdmins)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		} else {
			a.logger.Error("Error fetching user admins", zap.Error(err))
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userID": userAdmins.UserID,
		"admins": userAdmins.Admin,
	})
}
