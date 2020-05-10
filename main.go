package main

import (
	"github.com/Skamaniak/happiness-door-slack-bot/cmd"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/spf13/viper"
	"log"
)

func main() {
	log.Println("Starting app...")
	conf.InitConfig()
	port := viper.GetInt(conf.AppPort)
	cmd.Run(port)
}
