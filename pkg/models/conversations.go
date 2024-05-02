package models

import (
	"context"
	"database/sql"
	log "lexichat-backend/pkg/utils/logging"
	"time"
	// "github.com/google/uuid"
)

type MessageStatus string

const (
    Pending   MessageStatus = "pending"
    Sent      MessageStatus = "sent"
    Delivered MessageStatus = "delivered"
    Read      MessageStatus = "read"
)

type Message struct {
    ID           string `json:"id"`
    ChannelID    int64 	   `json:"channel_id"`
    CreatedAt    time.Time `json:"created_at"`
    SenderUserID string     `json:"sender_user_id"`
    Content      string    `json:"content"`
    Status       MessageStatus    `json:"status"`
}

func InsertMessage(message Message, db *sql.DB) error{
	errs := make(chan error, 1)
    go func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

		_, err := db.ExecContext(ctx, `
            INSERT INTO messages (sender_user_id, channel_id, content, id, status)
            VALUES ($1, $2, $3, $4)`,
            message.SenderUserID, message.ChannelID, message.Content, message.ID, Sent)
        if err != nil {
            log.ErrorLogger.Printf("Error inserting message into the database: %v", err)
			errs <- err
        }
    }()
	if err := <- errs; err != nil {
		return err
	}
	return nil
}

