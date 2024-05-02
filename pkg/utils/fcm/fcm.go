package utils

import (
	"context"
	"fmt"
	"lexichat-backend/pkg/models"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)


func SetupFCM() (*messaging.Client, error){
		opt := option.WithCredentialsFile("lexichat-ce3c1-firebase-adminsdk.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatalf("error while initializing FCM: %v\n", err)
			log.Panicf("error while initializing FCM: %v\n", err)
		}
	
		client, err := app.Messaging(context.Background())
		if err != nil {
			log.Panicf("error getting Messaging client: %v\n", err)
		}

		return client, err
}

func SendNotificationToClient(fcmToken string,client *messaging.Client ,title_message string, body string) {

    message := &messaging.Message{
        Notification: &messaging.Notification{
            Title: title_message,
            Body:  body,
        },
        // Token: "cBRJGauRQe-9GIqhKAx-Zh:APA91bGnSbRuFW-bIbJUSbYz75AnY-WhF7vntiu39i-L1uHCt_QrnEtoEkRi_qWoCLywugG8NwVgENYyYoqiggfi4E9yd-UETtbSEO1xTyrucl-aMmqgs-1avQ73l2QvyuIb-9AJdLbG", 
		Token: fcmToken,
    }

    response, err := client.Send(context.Background(), message)
    if err != nil {
        log.Fatalf("error sending message: %v\n", err)
    }

    fmt.Println("Successfully sent message:", response)
}

// func SendChannelDataToFcmClients(channelId int64, users []models.User, client *messaging.Client) {
// 	title := "####CREATE CHANNEL REQUEST####"
//     body := fmt.Sprintf("ChannelID: %d", channelId)
// 	for _, user := range(users) {
// 		SendMessage(user.FCMToken, client, title, body)
// 	}
// }

func SendNotificationsToClients(users []models.User, client *messaging.Client, title string, body string) {

	
	for _, user := range(users) {
		SendNotificationToClient(user.FCMToken, client, title, body)
	}
}


func messageToMap(msg models.Message) map[string]string {
	result := make(map[string]string)

	result["id"] = msg.ID
	result["channel_id"] = fmt.Sprintf("%d", msg.ChannelID)
	result["created_at"] = msg.CreatedAt.Format(time.RFC3339)
	result["sender_user_id"] = msg.SenderUserID
	result["content"] = msg.Content
	result["status"] = fmt.Sprintf("%v", msg.Status)

	return result
}



func SendDataAndNotificationsToClients(fcmTokens []string, client *messaging.Client, message models.Message, notificationTitle string, notificationBody string) error {
    
    if len(fcmTokens) == 0 {
        return nil
    }
    
    ctx := context.Background()

	
    // Create a message
    fcmMessage := &messaging.MulticastMessage{
        Tokens: fcmTokens,
        Data:   messageToMap(message),
        Notification: &messaging.Notification{
            Title: notificationTitle,
            Body:  notificationBody,
        },
    }

    batchResponse, err := client.SendMulticast(ctx, fcmMessage)
    if err != nil {
        return fmt.Errorf("failed to send FCM message: %w", err)
    }

    for _, response := range batchResponse.Responses {
        if response.Success {
            // Message was delivered successfully
            continue
			
        }

        // Handle error responses
        err := response.Error
        if err != nil {
            fmt.Printf("Failed to send message: %v\n", err)
        }
    }

    return nil
}

func SendDataToMultipleFCMClients(fcmTokens []string, client *messaging.Client, data map[string]string) error {
    // Create a new multicast message
    multicastMessage := &messaging.MulticastMessage{
        Tokens: fcmTokens,
        Data:   data,
    }

    // Send the multicast message with the client
    batchResponse, err := client.SendMulticast(context.Background(), multicastMessage)
    if err != nil {
        return err
    }

    // Check the batch response for success and failure counts
    if batchResponse.SuccessCount > 0 {
        fmt.Printf("Data sent to %d device(s)\n", batchResponse.SuccessCount)
    }

    if batchResponse.FailureCount > 0 {
        fmt.Printf("Failed to send data to %d device(s)\n", batchResponse.FailureCount)
        for _, resp := range batchResponse.Responses {
            if !resp.Success {
                fmt.Printf("Error: %s\n", resp.Error)
            }
        }
    }

    return nil
}