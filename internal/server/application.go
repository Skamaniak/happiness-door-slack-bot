package server

import (
	"github.com/Skamaniak/happiness-door-slack-bot/internal/api"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/service"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/ws"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func RegisterREST(r *mux.Router, s *service.SlackService) {
	p := viper.GetString(conf.RESTApiPrefix)
	logrus.WithField("prefix", p).Info("Registering REST API.")

	handlers := api.NewRestHandlers(s)
	r.HandleFunc(p+"/happiness-door", handlers.Initiation).
		Methods("POST")
	r.HandleFunc(p+"/happiness-door/interaction", handlers.Vote).
		Methods("POST")
}

func RegisterWS(r *mux.Router, s *service.SlackService) {
	p := viper.GetString(conf.WebApiPrefix)
	logrus.WithField("prefix", p).Info("Registering WEB socket API.")

	wsRouter := ws.NewRouter(s)
	handlers := api.NewWSHandlers(s)
	wsRouter.Handle(handlers.CreateVoteHandler())
	r.Handle(p+"/connect", wsRouter)
}

func RegisterWeb(r *mux.Router) {
	if viper.GetBool(conf.WebFileServerEnabled) {
		p := viper.GetString(conf.WebFileServerPrefix)
		logrus.WithField("prefix", p).Info("Registering WEB file server.")

		wf := viper.GetString(conf.WebFolder)
		fs := http.FileServer(http.Dir(wf))

		r.PathPrefix(p).Handler(http.StripPrefix(p, fs))
	}
}
