package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"lexichat-backend/pkg/models"
	"lexichat-backend/pkg/utils/logging"
	"net/http"
)

type _CreateChannelRequest struct {
    ChannelName  string 		`json:"channel_name"`
    TonalityTag  string 		`json:"tonality_tag"`
    Description  string 		`json:"description"`
	Users 		 []models.User  `json:"users"`

}

func CreateChannelHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var req _CreateChannelRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Insert into channels table
    var channelID int64
    channelInsertQuery := `INSERT INTO channels (channel_name, tonality_tag, description) VALUES ($1, $2, $3) RETURNING id`
    err = db.QueryRow(channelInsertQuery, req.ChannelName, req.TonalityTag, req.Description).Scan(&channelID)
    if err != nil {
        errorMessage := fmt.Sprintf("Failed to insert channel record: %v", err)
        http.Error(w, errorMessage, http.StatusInternalServerError)
        logging.ErrorLogger.Println(errorMessage)
        return
    }

	var userIDs []string

    // Insert into channel_users table
    for _, user := range req.Users {
        channelUserInsertQuery := `INSERT INTO channel_users (channel_id, user_id) VALUES ($1, $2)`
        _, err := db.Exec(channelUserInsertQuery, channelID, user.UserID)
        if err != nil {
            errorMessage := fmt.Sprintf("Failed to insert record into channel_users table: %v", err)
            http.Error(w, errorMessage, http.StatusInternalServerError)
            logging.ErrorLogger.Println(errorMessage)
            return
        }
		userIDs = append(userIDs, user.UserID)
    }

    // Prepare response
    response := map[string]interface{}{
        "channel_id": channelID,
        "user_ids":   userIDs,
    }
	jsonResponse, err := json.Marshal(response)
    if err != nil {
        errorMessage := fmt.Sprintf("Failed to marshal response while creating a channel: %v", err)
        http.Error(w, errorMessage, http.StatusInternalServerError)
        logging.ErrorLogger.Println(errorMessage)
        return
    }

	// send response(channel_id) to other users
	// use user's fcm to send responses.


    // Send response to creator user
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}