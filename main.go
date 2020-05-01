package main

import (
	"github.com/Skamaniak/happiness-door-slack-bot/cmd"
	"log"
	"os"
	"strconv"
)

const DefaultPort = 8080

func main() {
	var port int
	var err error

	envPortValue := os.Getenv("PORT")
	if envPortValue == "" {
		port = DefaultPort
	} else {
		port, err = strconv.Atoi(envPortValue)
		if err != nil {
			log.Println("Failed to parse port number", envPortValue, "using default", DefaultPort)
			port = DefaultPort
		}
	}

	log.Println("Starting app...")
	cmd.Run(port)
}
