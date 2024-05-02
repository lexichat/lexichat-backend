package utils

import (
	"database/sql"
	"fmt"
	"lexichat-backend/pkg/utils/logging"
)

func FetchFCMTokensOfUsers(db *sql.DB, userIDs []string) ([]string, error) {
	var fcmTokens []string

	for _, userID := range userIDs {
		var fcmToken string

		query := `SELECT fcm_token FROM users WHERE user_id = $1 `

        row := db.QueryRow(query, userID)

        err := row.Scan(&fcmToken)
        if err != nil {
			if err == sql.ErrNoRows {
                logging.ErrorLogger.Printf(fmt.Sprintf("No FCM token found for user ID: %s", userID))
            } else {
                logging.ErrorLogger.Printf(fmt.Sprintf("Failed to fetch FCM token for user ID %s: %v", userID, err))
            }
            continue 
        }

		fcmTokens = append(fcmTokens, fcmToken)

	}

    return fcmTokens, nil
}


func FetchUserIDsOfChannel(db *sql.DB, channelID int64) ([]string, error) {
    var userIDs []string

    query := `SELECT user_id FROM channel_users WHERE channel_id = $1`

    rows, err := db.Query(query, channelID)
    if err != nil {
        logging.ErrorLogger.Println("Failed to execute query:", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var userID string
        if err := rows.Scan(&userID); err != nil {
            logging.ErrorLogger.Println("Failed to scan row:", err)
            return nil, err
        }
        userIDs = append(userIDs, userID)
    }
    if err := rows.Err(); err != nil {
        logging.ErrorLogger.Println("Error during iteration:", err)
        return nil, err
    }

    return userIDs, nil
}