package utils

import (
	"context"
	"database/sql"
	"lexichat-backend/pkg/models"
	logging "lexichat-backend/pkg/utils/logging"
	// "github.com/google/uuid"
)

func StoreMessage(db *sql.DB, message models.Message) error {
	// message.ID, _ = uuid.NewRandom()
	message.Status = models.Sent

        _, err := db.ExecContext(context.Background(), `
            INSERT INTO messages (id, channel_id, sender_user_id, message, status)
            VALUES ($1, $2, $3, $4, $5)
        `, message.ID, message.ChannelID, message.SenderUserID, message.Content, message.Status)
        if err != nil {
            logging.ErrorLogger.Printf("Failed to store messages in db. %v", err)
        }     

    return err
}
