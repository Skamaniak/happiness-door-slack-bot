package conf

import "github.com/spf13/viper"

const AppPort = "PORT"
const SlackToken = "SLACK_TOKEN"
const DbHost = "DB_HOST"
const DbPort = "DB_PORT"
const DbName = "DB_NAME"
const DbUser = "DB_USER"
const DbPassword = "DB_PASSWORD"

func InitConfig() {
	viper.AutomaticEnv()

	// Slack
	viper.SetDefault(SlackToken, "")

	// App details
	viper.SetDefault(AppPort, 8080)

	// DB connection
	viper.SetDefault(DbHost, "localhost")
	viper.SetDefault(DbPort, 5432)
	viper.SetDefault(DbName, "happiness-door")
	viper.SetDefault(DbUser, "postgres")
	viper.SetDefault(DbPassword, "")
}
