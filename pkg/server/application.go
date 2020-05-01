package server

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/handler"
	"log"
	"net/http"
)

func RunServer(port int) {
	http.HandleFunc("/rest/v1/happiness-door", handler.HappinessDoorHandler)

	hostPort := fmt.Sprintf(":%d", port)
	log.Println("Registering handler to", hostPort)
	err := http.ListenAndServe(hostPort, nil)

	log.Println("ERR: Failed to create HTTP server", err)
}
