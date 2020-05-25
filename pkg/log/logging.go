package log

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func readLogLevelFromEnv() logrus.Level {
	lvl, ok := os.LookupEnv(conf.LogLevel)
	if !ok {
		lvl = "info"
	}
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.InfoLevel
	}
	return ll
}

func setLogLevelFromEnv() {
	ll := readLogLevelFromEnv()
	if logrus.GetLevel() != ll {
		logrus.WithFields(logrus.Fields{
			"OldLogLevel": logrus.GetLevel(),
			"NewLogLevel": ll,
		}).Info("Changing log level")
		logrus.SetLevel(ll)
	}
}

func InitLogging() {
	setLogLevelFromEnv()
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		logrus.Info("Starting periodic log level checker")
		for {
			select {
			case <-ticker.C:
				setLogLevelFromEnv()
			}
		}
	}()
}
