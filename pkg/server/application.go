package server

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/client"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/db"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/handler"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RunServer(port int) error {
	repo, err := db.NewHappinessDoor()
	if err != nil {
		return err
	}
	slackClient := client.NewSlackClient()
	slackService := service.NewSlackService(repo, slackClient)
	handlers := handler.NewHandlers(slackService)

	router := mux.NewRouter()
	router.HandleFunc("/rest/v1/happiness-door", handlers.Initiation).
		Methods("POST")
	router.HandleFunc("/rest/v1/happiness-door/interaction", handlers.Vote).
		Methods("POST")

	hostPort := fmt.Sprintf(":%d", port)
	logrus.WithField("Host", hostPort).Info("Registering handler")

	err = http.ListenAndServe(hostPort, router)
	return err
}
