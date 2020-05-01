package cmd

import "github.com/Skamaniak/happiness-door-slack-bot/pkg/server"

func Run(port int) {
	// Main HTTP server
	server.RunServer(port)

	//TODO admin server
}
