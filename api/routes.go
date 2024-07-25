package api

import (
	"net/http"

	"github.com/arjnep/goandchat/internal/chat"
	"github.com/arjnep/goandchat/internal/handlers"
)

func SetupRoutes(manager *chat.ChatroomManager, fs http.Handler) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(manager, w, r)
	})

	http.HandleFunc("/create_chatroom", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateChatroom(manager, w, r)
	})

	http.Handle("/", fs)
}
