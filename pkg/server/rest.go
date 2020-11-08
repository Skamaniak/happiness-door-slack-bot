package server

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/api"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func RunRestServer(s *service.SlackService) error {
	handlers := api.NewRestHandlers(s)

	router := mux.NewRouter()
	router.HandleFunc("/rest/v1/happiness-door", handlers.Initiation).
		Methods("POST")
	router.HandleFunc("/rest/v1/happiness-door/interaction", handlers.Vote).
		Methods("POST")

	hostPort := fmt.Sprintf(":%d", viper.GetInt(conf.RestPort))
	logrus.WithField("Host", hostPort).Info("Registering REST handler")

	return http.ListenAndServe(hostPort, router)
}
