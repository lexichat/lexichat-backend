package models

import (
	"database/sql"
	log "lexichat-backend/pkg/utils/logging"
)

type User struct {
    UserID        string	`json:"user_id"`
    Username      string    `json:"user_name"`
    PhoneNumber   string    `json:"phone_number"`
    ProfilePicture []byte   `json:"profile_picture"`
    CreatedAt     string    `json:"created_at"`
	FCMToken      string    `json:"fcm_token"`
}

func (u *User) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (user_id, user_name, phone_number, profile_picture, fcm_token) VALUES ($1, $2, $3, $4, $5)",
		u.UserID, u.Username, u.PhoneNumber, u.ProfilePicture, u.FCMToken)
	if err != nil {
		log.ErrorLogger.Println("Error in creating new user. ", err)
		return err
	}

	return nil
}