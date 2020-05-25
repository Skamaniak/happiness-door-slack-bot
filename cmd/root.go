package cmd

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/server"
	"github.com/sirupsen/logrus"
)

func Run(port int) {
	// Main HTTP server
	err := server.RunServer(port)

	if err != nil {
		logrus.WithError(err).Panic("Failed to initialise server")
	}
	//TODO admin server
}
