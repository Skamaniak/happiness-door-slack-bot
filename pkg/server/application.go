package server

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/db"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer(port int) error {
	repo, err := db.NewHappinessDoor()
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	handlers := handler.NewHandlers(repo)

	router.HandleFunc("/rest/v1/happiness-door", handlers.Initiation).
		Methods("POST")
	router.HandleFunc("/rest/v1/happiness-door/interaction", handlers.Vote).
		Methods("POST")

	hostPort := fmt.Sprintf(":%d", port)
	log.Println("Registering handler to", hostPort)

	err = http.ListenAndServe(hostPort, router)
	return err
}
