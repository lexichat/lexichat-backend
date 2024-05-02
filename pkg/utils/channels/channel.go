package utils

import (
	"database/sql"
	"fmt"

	"lexichat-backend/pkg/utils/logging"
)

func FetchChannelName(db *sql.DB, channelID int64) (string, error) {
    var channelName string

    query := `SELECT channel_name FROM channels WHERE id = $1`

    err := db.QueryRow(query, channelID).Scan(&channelName)
    if err != nil {
        if err == sql.ErrNoRows {
            logging.ErrorLogger.Println("No channel found for channel ID:", channelID)
            return "", fmt.Errorf("no channel name found for channel ID %d: %w", channelID, err)
        }
        logging.ErrorLogger.Println("Failed to fetch channel name:", err)
        return "", err
    }

    return channelName, nil
}
