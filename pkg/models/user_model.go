package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
    ID            uuid.UUID `json:"id"`
    Username      string    `json:"username"`
    PhoneNumber   string    `json:"phone_number"`
    ProfilePicture []byte   `json:"profile_picture"`
    CreatedAt     string    `json:"created_at"`
}

func (u *User) Create(db *sql.DB) error {
	u.ID = uuid.New()

	// Insert the user into the database
	_, err := db.Exec("INSERT INTO users (id, username, phone_number, profile_picture) VALUES ($1, $2, $3, $4)",
		u.ID, u.Username, u.PhoneNumber, u.ProfilePicture)
	if err != nil {
		return err
	}

	return nil
}