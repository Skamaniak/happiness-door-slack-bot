package main

import (
	"github.com/Skamaniak/happiness-door-slack-bot/cmd"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.Info("Starting app...")
	conf.InitConfig()
	log.InitLogging()

	port := viper.GetInt(conf.AppPort)
	cmd.Run(port)
}
