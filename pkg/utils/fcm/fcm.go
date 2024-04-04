package utils

import (
	"context"
	"fmt"
	"lexichat-backend/pkg/models"
	"log"

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

func SendMessage(fcmToken string,client *messaging.Client ,title_message string, body string) {

    message := &messaging.Message{
        Notification: &messaging.Notification{
            Title: title_message,
            Body:  body,
        },
        // Token: "dMEjAZ74TdSbGQAs-SwF8c:APA91bGcR9Hjw2ZACbjFiBNCMrj6bi-HNvzQq9HcUR4HWIbnZ3Nk5o_EIwJoGBVCzlskV6-cIwyr-P50EoE3t7jSeVWFhQfzQGHDB6ympNUC22Av42H3SITtD5zEwj7sf46bugq4dYeT", 
		Token: fcmToken,
    }

    response, err := client.Send(context.Background(), message)
    if err != nil {
        log.Fatalf("error sending message: %v\n", err)
    }

    fmt.Println("Successfully sent message:", response)
}

func SendChannelDataToFcmClients(channelId int64, users []models.User, client *messaging.Client) {
	title := "####CREATE CHANNEL REQUEST####"
    body := fmt.Sprintf("ChannelID: %d", channelId)
	for _, user := range(users) {
		SendMessage(user.FCMToken, client, title, body)
	}
}