package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	auth "lexichat-backend/pkg/utils/auth"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan Message
	ackChan  chan struct{}
	cancelFn context.CancelFunc
	UserID	 string
}

type Channel struct {
	ID      string
	clients map[*Client]bool
	broadcast chan Message
	register  chan *Client
	unregister chan *Client
	mu        sync.RWMutex
}

var channels = make(map[string]*Channel)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

type Acknowledgment struct {
	IsSent  bool  `json:"isSent"` 
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Sender   *Client 
}



func (c *Client) readPump(channel *Channel) {
	defer func() {
		c.cancelFn()
		channel.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error in reading: %v", err)
			return
		}
		msg.Sender = c
		channel.broadcast <- msg
		c.ackChan <- struct{}{}
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
			fmt.Println(channel.clients)

			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error in writing: %v", err)
				return
			}

		case <-c.ackChan:
			// Send acknowledgment here
			err := c.conn.WriteJSON(Acknowledgment{IsSent: true})
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
		}
	}
}

func (channel *Channel) runChannel() {
	for {
		select {
		case client := <-channel.register:
			channel.mu.Lock()
			channel.clients[client] = true
			channel.mu.Unlock()

			_, cancel := context.WithCancel(context.Background())
			client.cancelFn = cancel
			go client.writePump(channel)
			go client.readPump(channel)

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

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel")

	channel, ok := channels[channelID]
	if !ok {
		channel = &Channel{
			ID:        channelID,
			clients:   make(map[*Client]bool),
			broadcast: make(chan Message),
			register:  make(chan *Client),
			unregister: make(chan *Client),
		}
		channels[channelID] = channel
		go channel.runChannel()
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
		ackChan: make(chan struct{}, 1),
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
	channel, ok := channels[channelID]
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
