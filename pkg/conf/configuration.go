package conf

import (
	"github.com/spf13/viper"
)

const AppPort = "PORT"
const SlackToken = "SLACK_TOKEN"
const WebTokenLength = "WEB_TOKEN_LENGTH"
const DbUrl = "DATABASE_URL"
const LogLevel = "LOG_LEVEL"

func InitConfig() {
	viper.AutomaticEnv()

	// Slack
	viper.SetDefault(SlackToken, "")

	// App details
	viper.SetDefault(AppPort, 8080)

	// Access from web
	viper.SetDefault(WebTokenLength, 128)

	// DB connection
	viper.SetDefault(DbUrl, "postgres://postgres:@localhost:5432/happiness-door")
}
