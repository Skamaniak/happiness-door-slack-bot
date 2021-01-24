package conf

import (
	"github.com/spf13/viper"
	"time"
)

const (
	Port                      = "PORT"
	SlackToken                = "SLACK_TOKEN"
	DbUrl                     = "DATABASE_URL"
	LogLevel                  = "LOG_LEVEL"
	BotName                   = "BOT_NAME"
	RESTApiPrefix             = "REST_API_PREFIX"
	WebTokenLength            = "WEB_TOKEN_LENGTH"
	WebHost                   = "WEB_HOST"
	WebScheme                 = "WEB_SCHEME"
	WebFolder                 = "WEB_FOLDER"
	WebFileServerEnabled      = "WEB_FILE_SERVER_ENABLED"
	WebFileServerPrefix       = "WEB_FILE_SERVER_PREFIX"
	WebApiPrefix              = "WEB_API_PREFIX"
	WebSocketPingPongInterval = "WEB_PING_PONG_INTERVAL"
	WebSocketMaxPingPongDelay = "WEB_PING_PONG_MAX_DELAY"
	AdminApiEnabled           = "ADMIN_API_ENABLED"
	AdminPort                 = "ADMIN_PORT"
	AdminApiPrefix            = "ADMIN_API_PREFIX"
)

func InitConfig() {
	viper.AutomaticEnv()

	// Bot details
	viper.SetDefault(BotName, "happiness_door_bot")

	// Slack
	viper.SetDefault(SlackToken, "")

	// App details
	viper.SetDefault(Port, 8080)

	// REST
	viper.SetDefault(RESTApiPrefix, "/rest/v1")

	// Web
	viper.SetDefault(WebTokenLength, 128)
	viper.SetDefault(WebHost, "localhost:"+viper.GetString(Port))
	viper.SetDefault(WebScheme, "http")
	viper.SetDefault(WebApiPrefix, "/ws/v1")
	viper.SetDefault(WebSocketPingPongInterval, 30*time.Second)
	viper.SetDefault(WebSocketMaxPingPongDelay, viper.GetDuration(WebSocketPingPongInterval)*3)

	viper.SetDefault(WebFileServerEnabled, true)
	viper.SetDefault(WebFolder, "./frontend/dist")
	viper.SetDefault(WebFileServerPrefix, "/")

	// Admin
	viper.SetDefault(AdminPort, 8079)
	viper.SetDefault(AdminApiEnabled, false)
	viper.SetDefault(AdminApiPrefix, "/admin")

	// DB connection
	viper.SetDefault(DbUrl, "postgres://postgres:@localhost:5432/happiness-door")
}
