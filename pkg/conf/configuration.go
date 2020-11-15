package conf

import (
	"github.com/spf13/viper"
)

const Port = "PORT"
const SlackToken = "SLACK_TOKEN"
const WebTokenLength = "WEB_TOKEN_LENGTH"
const DbUrl = "DATABASE_URL"
const LogLevel = "LOG_LEVEL"
const BotName = "BOT_NAME"

func InitConfig() {
	viper.AutomaticEnv()

	// Bot details
	viper.SetDefault(BotName, "happiness_door_bot")

	// Slack
	viper.SetDefault(SlackToken, "")

	// App details
	viper.SetDefault(Port, 8080)

	// Web auth
	viper.SetDefault(WebTokenLength, 128)

	// DB connection
	viper.SetDefault(DbUrl, "postgres://postgres:@localhost:5432/happiness-door")
}
