package cmd

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/client"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/db"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/server"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func createSlackService() (*service.SlackService, error) {
	repo, err := db.NewHappinessDoor()
	if err != nil {
		return nil, err
	}
	slackClient := client.NewSlackClient()
	return service.NewSlackService(repo, slackClient), nil
}

func startServer(r *mux.Router, hostPort string) {
	go func() {
		logrus.WithField("Host", hostPort).Info("Registering HTTP handler.")
		err := http.ListenAndServe(hostPort, r)
		if err != nil {
			logrus.WithError(err).Panic("Failed to initialise server.")
		}
	}()
}

func awaitTermination() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func Run() {
	s, err := createSlackService()
	if err != nil {
		logrus.WithError(err).Panic("Failed to create Slack service.")
	}

	r := mux.NewRouter()
	server.RegisterREST(r, s)
	server.RegisterWS(r, s)
	server.RegisterWeb(r)
	hp := fmt.Sprintf(":%d", viper.GetInt(conf.Port))
	startServer(r, hp)

	if viper.GetBool(conf.AdminApiEnabled) {
		r = mux.NewRouter()
		server.RegisterAdminAPI(r, s)
		hp = fmt.Sprintf(":%d", viper.GetInt(conf.AdminPort))
		startServer(r, hp)
	}

	awaitTermination()
	logrus.Info("Stopping Happiness Door Bot.")
}
