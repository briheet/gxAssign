package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/briheet/gxAssign/internal/models"
	"github.com/briheet/gxAssign/internal/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		a.logger.Error("adminEmail or adminPassword issue in register", zap.Error(err))

		http.Error(w, "Email or Password not found", http.StatusBadRequest)
		return
	}

	filter := bson.M{"email": admin.Email}
	var existingAdmin models.Admin

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
		a.logger.Error("Error creating a admin", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "admin created successfully"})
}

func (a *api) adminLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var adminDetails models.AdminAuth
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
	var existingAdminDetails *models.Admin

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

func (a *api) acceptAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	assignmentID := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(assignmentID)
	if err != nil {
		a.logger.Error("invalid ObjectId", zap.Error(err))
		http.Error(w, "invalid assignment ID", http.StatusBadRequest)
		return
	}

	var admin models.Admin
	err = json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		a.logger.Error("error decoding admin details in acceptAssignments", zap.Error(err))
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if admin.Email == "" {
		a.logger.Error("admin email is empty", zap.Error(err))
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	oldCollection := a.database.Collection("documents")
	filter := bson.M{"admin": admin.Email, "_id": objectID}

	var userDocuments []models.UserDocument
	cursor, err := oldCollection.Find(ctx, filter)
	if err != nil {
		a.logger.Error("error finding user documents", zap.Error(err))
		http.Error(w, "error fetching user documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &userDocuments); err != nil {
		a.logger.Error("error decoding user documents", zap.Error(err))
		http.Error(w, "error decoding user documents", http.StatusInternalServerError)
		return
	}

	if len(userDocuments) == 0 {
		http.Error(w, "no assignments found for this admin", http.StatusNotFound)
		return
	}

	newCollection := a.database.Collection("accepted_documents")
	acceptedDocuments := make([]interface{}, len(userDocuments))
	for i, doc := range userDocuments {
		acceptedDocuments[i] = bson.M{
			"userID":     doc.UserID,
			"task":       doc.Task,
			"admin":      doc.Admin,
			"status":     "accepted",
			"acceptedAt": time.Now(),
		}
	}

	_, err = newCollection.InsertMany(ctx, acceptedDocuments)
	if err != nil {
		a.logger.Error("error inserting into accepted_documents", zap.Error(err))
		http.Error(w, "error processing user documents", http.StatusInternalServerError)
		return
	}

	_, err = oldCollection.DeleteMany(ctx, filter)
	if err != nil {
		a.logger.Error("error deleting user documents from documents", zap.Error(err))
		http.Error(w, "error deleting user documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "assignments accepted and moved successfully"})
}

func (a *api) rejectAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	assignmentID := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(assignmentID)
	if err != nil {
		a.logger.Error("invalid ObjectId", zap.Error(err))
		http.Error(w, "invalid assignment ID", http.StatusBadRequest)
		return
	}

	var admin models.Admin
	err = json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		a.logger.Error("error decoding admin details in rejectAssignments", zap.Error(err))
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if admin.Email == "" {
		a.logger.Error("admin email is empty", zap.Error(err))
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	oldCollection := a.database.Collection("documents")
	filter := bson.M{"admin": admin.Email, "_id": objectID}

	var userDocuments []models.UserDocument
	cursor, err := oldCollection.Find(ctx, filter)
	if err != nil {
		a.logger.Error("error finding user documents", zap.Error(err))
		http.Error(w, "error fetching user documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &userDocuments); err != nil {
		a.logger.Error("error decoding user documents", zap.Error(err))
		http.Error(w, "error decoding user documents", http.StatusInternalServerError)
		return
	}

	if len(userDocuments) == 0 {
		http.Error(w, "no assignments found for this admin", http.StatusNotFound)
		return
	}

	rejectedCollection := a.database.Collection("rejected_documents")
	rejectedDocuments := make([]interface{}, len(userDocuments))
	for i, doc := range userDocuments {
		rejectedDocuments[i] = bson.M{
			"userID":          doc.UserID,
			"task":            doc.Task,
			"admin":           doc.Admin,
			"status":          "rejected",
			"rejectedAt":      time.Now(),
			"rejectionReason": "Assignment rejected by admin",
		}
	}

	_, err = rejectedCollection.InsertMany(ctx, rejectedDocuments)
	if err != nil {
		a.logger.Error("error inserting into rejected_documents", zap.Error(err))
		http.Error(w, "error processing user documents", http.StatusInternalServerError)
		return
	}

	_, err = oldCollection.DeleteMany(ctx, filter)
	if err != nil {
		a.logger.Error("error deleting user documents from documents", zap.Error(err))
		http.Error(w, "error deleting user documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "assignments rejected and moved successfully"})
}
