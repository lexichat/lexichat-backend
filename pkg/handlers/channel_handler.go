package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"lexichat-backend/pkg/models"
	utils "lexichat-backend/pkg/utils/fcm"
	"lexichat-backend/pkg/utils/logging"
	"net/http"

	"firebase.google.com/go/messaging"
)

type _CreateChannelRequest struct {
    ChannelName         string 		    `json:"channel_name"`
    TonalityTag         string 		    `json:"tonality_tag"`
    Description         string 		    `json:"description"`
	Users 		        []string       `json:"users"`
    SenderUserID        string         `json:"sender_user_id"`
    FirstMessageContent string         `json:"first_message"`
    MessageID           string         `json:"message_id"`
}

func CreateChannelHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, fcmClient *messaging.Client) {
    var req _CreateChannelRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    fmt.Sprintf("req: %#v", req)

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
        _, err := db.Exec(channelUserInsertQuery, channelID, user)
        if err != nil {
            errorMessage := fmt.Sprintf("Failed to insert record into channel_users table: %v", err)
            http.Error(w, errorMessage, http.StatusInternalServerError)
            logging.ErrorLogger.Println(errorMessage)
            return
        }
		userIDs = append(userIDs, user)
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
    channel_data := models.Channel{
        ID:          channelID,
        Name:        req.ChannelName,
        CreatedAt:   time.Now(),
        TonalityTag: req.TonalityTag,
        Description: req.Description,
    }

    var sender_user_data models.User
    fmt.Println("sender user id: ", req.SenderUserID)
    sender_user_data, _ = FetchUserDetails(req.SenderUserID, db)
    sendNewChannelRequestToUsers(db, channel_data, sender_user_data, fcmClient, userIDs, models.Message{ID: req.MessageID, ChannelID: channelID, CreatedAt: time.Now(), SenderUserID: req.SenderUserID, Content: req.FirstMessageContent, Status: models.Sent})  


    // Send response to creator user
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}

func sendNewChannelRequestToUsers(db *sql.DB, channel_data models.Channel, sender_user_data models.User, fcmClient *messaging.Client, userIDs []string, firstMessage models.Message) {
    // Construct data to send
	channelDataJSON, _ := json.Marshal(channel_data)
	senderUserDataJSON, _ := json.Marshal(sender_user_data)
    firstMessageData, _ := json.Marshal(firstMessage)

	flattenedUserIDs := strings.Join(userIDs, ",")
	data := map[string]string{
		"NewUserChat": fmt.Sprintf(`{"channel": %s, "sender_user": %s, "channel_users": "%s", "first_message": %s}`, string(channelDataJSON), string(senderUserDataJSON), flattenedUserIDs, string(firstMessageData)),
	}

    fmt.Println(data);

    var fcmTokens []string
    for _, userID := range userIDs {
        var fcmToken string
        fetchFCMTokensQuery := `SELECT fcm_token FROM users WHERE user_id = $1`
        err := db.QueryRow(fetchFCMTokensQuery, userID).Scan(&fcmToken)
        if err != nil {
            // Handle error
            fmt.Printf("Failed to fetch FCM token for user %s: %v\n", userID, err)
            continue
        }
        if (len(fcmToken) != 0) {
            fcmTokens = append(fcmTokens, fcmToken)
        }
    }


    err :=  utils.SendDataToMultipleFCMClients(fcmTokens, fcmClient, data)
    if err != nil {
        // Handle error
        fmt.Printf("Failed to send data to FCM clients: %v\n", err)
    }
}