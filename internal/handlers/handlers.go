package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/arjnep/goandchat/internal/chat"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type CreateChatroomResponse struct {
	Code string `json:"code"`
}

func ServeWs(manager *chat.ChatroomManager, w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	hub, exists := manager.GetChatroom(code)
	if !exists {
		http.Error(w, "Chatroom not found", http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	client := chat.NewClient(hub, conn, code)
	hub.Register <- client
	go client.WriteMessages()
	go client.ReadMessages()
}

func CreateChatroom(manager *chat.ChatroomManager, w http.ResponseWriter, r *http.Request) {
	code := generateUniqueCode()
	manager.CreateChatroom(code)
	response := CreateChatroomResponse{Code: code}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateUniqueCode() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	code := make([]rune, 5)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
