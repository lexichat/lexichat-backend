package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"lexichat-backend/pkg/models"
	jwt "lexichat-backend/pkg/utils/auth"
	log "lexichat-backend/pkg/utils/logging"
)


func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.ErrorLogger.Printf("Error decoding JSON: %v\n", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = user.Create(db)
		if err != nil {
			log.ErrorLogger.Printf("Error creating user: %v\n", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		token, err := jwt.GenerateJWT(user.ID.String())
		if err != nil {
			log.ErrorLogger.Printf("Error generating JWT token: %v\n", err)
			http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
			return
		}

		log.InfoLogger.Printf("User created successfully: %s\n", user.ID)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"userId": user.ID.String(),
			"token":  token,
		})
	}
}