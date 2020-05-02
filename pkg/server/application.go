package server

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/rest/v1/happiness-door", handler.HappinessDoorHandler).
		Methods("POST")

	hostPort := fmt.Sprintf(":%d", port)
	log.Println("Registering handler to", hostPort)
	err := http.ListenAndServe(hostPort, router)

	log.Println("ERR: Failed to create HTTP server", err)
}
