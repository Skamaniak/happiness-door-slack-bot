package conf

import "github.com/spf13/viper"

const AppPort = "PORT"
const SlackToken = "SLACK_TOKEN"
const DbUrl = "DATABASE_URL"

func InitConfig() {
	viper.AutomaticEnv()

	// Slack
	viper.SetDefault(SlackToken, "")

	// App details
	viper.SetDefault(AppPort, 8080)

	// DB connection
	viper.SetDefault(DbUrl, "postgres://postgres:@localhost:5432/happiness-door")
}
