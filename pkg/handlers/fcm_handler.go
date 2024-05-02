package handlers

import (
	"database/sql"
	"fmt"
	"lexichat-backend/pkg/models"
	channel_util "lexichat-backend/pkg/utils/channels"
	fcm "lexichat-backend/pkg/utils/fcm"
	user "lexichat-backend/pkg/utils/users"
	"sync"

	"firebase.google.com/go/messaging"
)

type NotificationData struct {
    Title string 
    Body  string 
}


func OffLoadMessageToFCM(fcmClient *messaging.Client, message models.Message, dbClient *sql.DB, channel *Channel) {
	var wg sync.WaitGroup
	var channelName string
	var userFCMTokens []string
	var err error
	var notification NotificationData


	wg.Add(2)

	go func() {
        defer wg.Done()
        channelName, err = channel_util.FetchChannelName(dbClient, message.ChannelID)
        if err != nil {
            fmt.Printf("Error getting channel name: %v\n", err)
            return
        }
    }()

    go func() {
        defer wg.Done()
        userFCMTokens, err = fetchUnconnectedClientFCMs(dbClient, channel)
        if err != nil {
            fmt.Printf("Error fetching user tokens: %v\n", err)
            return
        }
    }()

    wg.Wait()

	if err != nil {
		return
	}
	
	notification.Title = channelName
	notification.Body = message.Content


	err = fcm.SendDataAndNotificationsToClients(userFCMTokens, fcmClient, message, notification.Title, notification.Body)
	if err != nil {
		return
	}

}


func fetchUnconnectedClientFCMs(dbClient *sql.DB, channel *Channel) ([]string, error){
	unconnectedUserIds, err := channel.FetchUnConnectedClientsIds(dbClient)
	if err != nil {
		fmt.Printf("failed to fetch unConnected Clients. %v", err)
		return nil,  err
	}

	fcmTokens, err := user.FetchFCMTokensOfUsers(dbClient, unconnectedUserIds)
	if err != nil {
		fmt.Printf("failed to fetch fcm tokens of userIDs. %v", err)
		return nil, err
	}
	return fcmTokens, nil
}





	// store message in db

	// send back ack with status change

	// get channel name from channel id 

	// fetch users[fcm_token] from user channel table

	// remove senderUserId from users

	// send notifications 
	

