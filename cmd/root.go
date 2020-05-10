package cmd

import "github.com/Skamaniak/happiness-door-slack-bot/pkg/server"

func Run(port int) {
	// Main HTTP server
	err := server.RunServer(port)

	if err != nil {
		panic(err)
	}
	//TODO admin server
}
