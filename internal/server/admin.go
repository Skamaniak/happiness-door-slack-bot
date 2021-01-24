package server

import (
	"github.com/Skamaniak/happiness-door-slack-bot/internal/api"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterAdminAPI(r *mux.Router, s *service.SlackService) {
	p := viper.GetString(conf.AdminApiPrefix)
	logrus.WithField("prefix", p).Info("Registering ping API.")

	handlers := api.NewAdminHandlers(s)
	r.HandleFunc(p+"/ping", handlers.Ping).
		Methods("GET")
	r.HandleFunc(p+"/health", handlers.Health).
		Methods("GET")
}
