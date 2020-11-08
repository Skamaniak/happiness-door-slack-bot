package server

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/api"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/ws"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func RunWsServer(s *service.SlackService) error {
	router := ws.NewRouter(s)

	handlers := api.NewWSHandlers(s)

	router.Handle(handlers.CreateVoteHandler())

	http.Handle("/", router)

	hostPort := fmt.Sprintf(":%d", viper.GetInt(conf.WsPort))
	logrus.WithField("Host", hostPort).Info("Registering WebSocket handler")
	return http.ListenAndServe(hostPort, nil)
}
