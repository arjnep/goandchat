package main

import (
	"log"
	"net/http"
	"os"

	"github.com/arjnep/goandchat/api"
	"github.com/arjnep/goandchat/internal/chat"
)

func main() {
	manager := chat.NewChatroomManager()

	api.SetupRoutes(manager)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server started on", ":"+port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
