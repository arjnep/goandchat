package chat

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Messages   []Message
	Mutex      sync.Mutex
	UserCount  int
	Manager    *ChatroomManager
}

func NewHub(m *ChatroomManager) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Messages:   []Message{},
		UserCount:  0,
		Manager:    m,
	}
}
func (h *Hub) Run(code string) {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client, code)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) cleanup(code string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Manager.chatrooms, code)
	log.Println("Chatroom Deleted!")
}

func (h *Hub) registerClient(client *Client) {
	h.Clients[client] = true
	h.UserCount++
	log.Printf("Client registered. Total clients: %d", h.UserCount)
	h.broadcastUserCount()
	h.sendPreviousMessages(client)
}

func (h *Hub) unregisterClient(client *Client, code string) {
	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
		close(client.Send)
		h.UserCount--
		log.Printf("Client unregistered. Total clients: %d", h.UserCount)
		if len(h.Clients) == 0 {
			h.cleanup(code)
		}
		h.broadcastUserCount()
	}
}

func (h *Hub) broadcastMessage(message Message) {
	h.storeMessage(message)
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
			h.UserCount--
			log.Printf("Client unregistered due to send failure. Total clients: %d", h.UserCount)
			h.broadcastUserCount()
		}
	}
}

func (h *Hub) broadcastUserCount() {
	userCountMsg := Message{Sender: "system", Content: fmt.Sprintf("User count: %d", h.UserCount), Time: time.Now().Unix()}
	for client := range h.Clients {
		select {
		case client.Send <- userCountMsg:
		default:
			close(client.Send)
			delete(h.Clients, client)
			h.UserCount--
			log.Printf("Client unregistered due to send failure 2. Total clients: %d", h.UserCount)
		}
	}
}

type ChatroomManager struct {
	chatrooms map[string]*Hub
	mutex     sync.Mutex
}

func NewChatroomManager() *ChatroomManager {
	return &ChatroomManager{
		chatrooms: make(map[string]*Hub),
	}
}

func (m *ChatroomManager) CreateChatroom(code string) *Hub {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	hub := NewHub(m)
	m.chatrooms[code] = hub
	go hub.Run(code)
	return hub
}

func (m *ChatroomManager) GetChatroom(code string) (*Hub, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	hub, exists := m.chatrooms[code]
	return hub, exists
}

func (h *Hub) storeMessage(message Message) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Messages = append(h.Messages, message)
}

func (h *Hub) sendPreviousMessages(client *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	for _, message := range h.Messages {
		client.Send <- message
	}
}
