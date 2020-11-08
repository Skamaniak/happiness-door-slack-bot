package cmd

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/client"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/db"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/server"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/sirupsen/logrus"
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

func initServer(runServer func() error) {
	go func() {
		err := runServer()
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

	// Run Web Socket for frontend communication
	initServer(func() error {
		return server.RunWsServer(s)
	})

	// Main HTTP server
	initServer(func() error {
		return server.RunRestServer(s)
	})

	//TODO admin server

	awaitTermination()
	logrus.Info("Stopping Happiness Door Bot.")
}
