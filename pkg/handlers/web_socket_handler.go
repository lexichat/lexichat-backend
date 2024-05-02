package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"lexichat-backend/pkg/models"
	auth "lexichat-backend/pkg/utils/auth"
	chats "lexichat-backend/pkg/utils/chats"
	"log"
	"net/http"
	"sync"

	"firebase.google.com/go/messaging"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan Message
	ackChan  chan Acknowledgment
	cancelFn context.CancelFunc
	UserID	 string
}

type Channel struct {
	ID      int64
	clients map[*Client]bool
	broadcast chan Message
	register  chan *Client
	unregister chan *Client
	mu        sync.RWMutex
}

var channels = make(map[int64]*Channel)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

type Acknowledgment struct {
	Status  	models.MessageStatus  `json:"status"` 
	MessageID 	string 				  `json:"message_id"`
}

type Message struct {
	Sender   *Client 
	Message  InBoundMessage
}

type InBoundMessage struct {
	Content string `json:"message"`
	MessageID string `json:"message_id"`
}

type OutBoundMessage struct {

}



func (c *Client) readPump(channel *Channel, dbClient *sql.DB, fcmClient *messaging.Client) {
	defer func() {
		c.cancelFn()
		channel.unregister <- c
		c.conn.Close()
	}()

	for {
		var inBoundMsg InBoundMessage
		var msg Message
		var ack Acknowledgment
		err := c.conn.ReadJSON(&inBoundMsg)
		if err != nil {
			log.Printf("error in reading: %v", err)
			return
		}
		
		msg.Message = inBoundMsg
		msg.Sender = c
		channel.broadcast <- msg

		fmt.Println(msg.Message)

    	// store msg in db
		var dbMsg models.Message
		dbMsg.ChannelID = channel.ID
		dbMsg.ID = msg.Message.MessageID
		dbMsg.Content = msg.Message.Content
		dbMsg.SenderUserID = msg.Sender.UserID

		var wg sync.WaitGroup

		wg.Add(2)

		go func ()  {
			chats.StoreMessage(dbClient, dbMsg)
			defer wg.Done()
		}()		

		ack.MessageID = dbMsg.ID
        c.ackChan <- ack


		go func ()  {
			// OffLoadMessageToFCM(fcmClient, dbMsg, dbClient, channel)
			defer wg.Done()
		}()
		wg.Wait()
	}
}

func (c *Client) writePump(channel *Channel) {
	defer c.conn.Close()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}

			var outboundmsg models.Message
			outboundmsg.ID = msg.Message.MessageID
			outboundmsg.SenderUserID = msg.Sender.UserID
			outboundmsg.CreatedAt = time.Now()
			outboundmsg.Status = models.Sent
			outboundmsg.ChannelID = channel.ID
			outboundmsg.Content = msg.Message.Content

			outboundData := map[string]interface{}{
				"Message" : outboundmsg,
			}

			err := c.conn.WriteJSON(outboundData)
			if err != nil {
				log.Printf("error in writing: %v", err)
				return
			}

		case ack := <-c.ackChan:
			sentStatusUpdate := map[string]interface{} {
				"Status" : Acknowledgment{Status: models.Sent, MessageID: ack.MessageID},
			}
			err := c.conn.WriteJSON(sentStatusUpdate)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
		}
	}
}

func (channel *Channel) runChannel(dbClient *sql.DB, fcmClient *messaging.Client) {
	for {
		select {
		case client := <-channel.register:
			channel.mu.Lock()
			channel.clients[client] = true
			channel.mu.Unlock()

			_, cancel := context.WithCancel(context.Background())
			client.cancelFn = cancel
			go client.writePump(channel)
			go client.readPump(channel, dbClient, fcmClient)

		case client := <-channel.unregister:
			channel.mu.Lock()
			if _, ok := channel.clients[client]; ok {
				delete(channel.clients, client)
				close(client.send)
			}
			channel.mu.Unlock()

		case msg := <-channel.broadcast:
			channel.mu.RLock()
			for client := range channel.clients {

				if client.UserID != msg.Sender.UserID {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(channel.clients, client)
				
			}
				}}
			channel.mu.RUnlock()
		}
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request, dbClient *sql.DB, fcmClient *messaging.Client) {
	channelID := r.URL.Query().Get("channel")
	channelIDInt, _ := strconv.ParseInt(channelID, 10, 64)

	channel, ok := channels[channelIDInt]
	if !ok {
		channel = &Channel{
			ID:        channelIDInt,
			clients:   make(map[*Client]bool),
			broadcast: make(chan Message),
			register:  make(chan *Client),
			unregister: make(chan *Client),
		}
		channels[channelIDInt] = channel
		go channel.runChannel(dbClient, fcmClient)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	tokenString := r.Header.Get("Authorization")
	userId, err := auth.GetUserIdFromToken(tokenString)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	client := &Client{
		conn:    conn,
		send:    make(chan Message, 256),
		ackChan: make(chan Acknowledgment, 8),
		UserID: userId,
	}

	channel.register <- client
}

func (channel *Channel) GetActiveClients() []*Client {
	channel.mu.RLock()
	defer channel.mu.RUnlock()

	clients := make([]*Client, 0, len(channel.clients))
	for client := range channel.clients {
		clients = append(clients, client)
	}
	return clients
}

func GetActiveClientsHandler(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel")
	channelIDInt, _ := strconv.ParseInt(channelID, 10, 64)
	channel, ok := channels[channelIDInt]
	if !ok {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	activeClients := channel.GetActiveClients()

	response := struct {
		Clients []*Client `json:"clients"`
	}{
		Clients: activeClients,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func (channel *Channel) FetchUnConnectedClientsIds(dbClient *sql.DB) ([]string, error) {
    ctx := context.Background()

    rows, err := dbClient.QueryContext(ctx, `
        SELECT user_id
        FROM channel_users
        WHERE channel_id = $1
    `, channel.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch unconnected clients: %w", err)
    }
    defer rows.Close()

    var dbClients []string
    for rows.Next() {
        var userID string
        err = rows.Scan(&userID)
        if err != nil {
            return nil, fmt.Errorf("failed to scan user ID: %w", err)
        }
        dbClients = append(dbClients, userID)
    }

    activeClients := channel.GetActiveClients()

    var unconnectedUserIds []string
    for _, dbclient := range dbClients {
        found := false
        for _, activeClient := range activeClients {
            if activeClient.UserID == dbclient {
                found = true
                break
            }
        }
        if !found {
            unconnectedUserIds = append(unconnectedUserIds, dbclient)
        }
    }

    return unconnectedUserIds, nil
}
