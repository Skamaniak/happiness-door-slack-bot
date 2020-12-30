package main

import (
	"github.com/Skamaniak/happiness-door-slack-bot/cmd"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/log"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting app...")
	conf.InitConfig()
	log.InitLogging()

	cmd.Run()
}
