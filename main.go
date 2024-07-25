package main

import (
	"log"
	"net/http"
	"os"
	"embed"
	"io/fs"
	"github.com/arjnep/goandchat/api"
	"github.com/arjnep/goandchat/internal/chat"
)

//go:embed static
var static embed.FS

func main() {
	var err error
	var staticFS = fs.FS(static)

	content, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal("Embeding Static: ", err)
	}
	fs := http.FileServer(http.FS(content))

	manager := chat.NewChatroomManager()

	api.SetupRoutes(manager, fs)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server started on", ":"+port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
