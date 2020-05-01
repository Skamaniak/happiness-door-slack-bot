package cmd

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/server"
	"log"
	"net/http"
)

func Run(port int) {
	http.HandleFunc("/", server.HelloServer)

	hostPort := fmt.Sprintf(":%d", port)
	log.Println("Registering handler to", hostPort)
	err := http.ListenAndServe(hostPort, nil)

	log.Println("ERR: Failed to create HTTP server", err)
}
