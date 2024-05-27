package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan Message
}

func NewClient(hub *Hub, conn *websocket.Conn, roomCode string) *Client {
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan Message),
	}

	go func() {
		client.Send <- Message{
			Sender:  "system",
			Content: "Room code: " + roomCode,
			Time:    time.Now().Unix(),
		}
	}()

	return client
}

func (c *Client) ReadMessages() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}
		msg.Time = time.Now().Unix()
		c.Hub.Broadcast <- msg
	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.Conn.Close()
	}()
	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("error writing message: %v", err)
			break
		}
	}
}
