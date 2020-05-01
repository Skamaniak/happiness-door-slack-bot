package cmd

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/server"
	"log"
	"net/http"
)

func Run() {
	http.HandleFunc("/", server.HelloServer)
	err := http.ListenAndServe(":8080", nil)

	log.Fatal("Failed to create HTTP server", err)
}
