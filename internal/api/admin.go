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

func (a *api) adminRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var admin *models.Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		a.logger.Error("error via decoding Admin details in register", zap.Error(err))
		http.Error(w, "invaild request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if admin.Email == "" || admin.Password == "" {
		a.logger.Error("userEmail or userPassword issue in register", zap.Error(err))
		http.Error(w, "Email or Password not found", http.StatusBadRequest)
		return
	}

	filter := bson.M{"email": admin.Email}
	var existingAdmin models.User

	err = a.database.Collection("admins").FindOne(ctx, filter).Decode(&existingAdmin)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Error("error checking for admin in register", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if existingAdmin.Email != "" {
		http.Error(w, "Admin already exists", http.StatusConflict)
		return

	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		a.logger.Error("error hashing password", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	admin.Password = hashedPassword

	_, err = a.database.Collection("admins").InsertOne(ctx, admin)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		a.logger.Error("Error creating a user", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
}

func (a *api) adminLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var adminDetails models.UserAuth
	err := json.NewDecoder(r.Body).Decode(&adminDetails)
	if err != nil {
		a.logger.Error("error via decoding admin details in register", zap.Error(err))
		http.Error(w, "invaild request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if adminDetails.Email == "" || adminDetails.Password == "" {
		http.Error(w, "Email or Password not found", http.StatusBadRequest)
		return
	}

	filter := bson.M{"email": adminDetails.Email}
	var existingAdminDetails *models.User

	err = a.database.Collection("admins").FindOne(ctx, filter).Decode(&existingAdminDetails)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
		} else {
			a.logger.Error("error finding admin in login", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	if !utils.CheckPasswordHash(adminDetails.Password, existingAdminDetails.Password) {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
}

func (a *api) assignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var admin models.Admin

	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		a.logger.Error("error via decoding admin details in register", zap.Error(err))
		http.Error(w, "invaild request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if admin.Email == "" {
		a.logger.Error("admin email is empty", zap.Error(err))
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	collection := a.database.Collection("documents")
	filter := bson.M{"email": admin.Email}

	var assignments []models.UserDocument
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		a.logger.Error("error finding assignments", zap.Error(err))
		http.Error(w, "error fetching assignments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &assignments); err != nil {
		a.logger.Error("error decoding assignments", zap.Error(err))
		http.Error(w, "error decoding assignments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(assignments)
	if err != nil {
		a.logger.Error("error encoding response", zap.Error(err))
		http.Error(w, "error generating response", http.StatusInternalServerError)
		return
	}
}
