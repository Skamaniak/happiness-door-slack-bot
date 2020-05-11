package client

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type SlackClient struct {
	api *slack.Client
}

func NewSlackClient() *SlackClient {
	token := viper.GetString(conf.SlackToken)
	api := slack.New(token)
	return &SlackClient{api: api}
}

func (c *SlackClient) GetUserIconUrl(userId string) (string, error) {
	profile, err := c.api.GetUserProfile(userId, false)
	if err != nil {
		return "", err
	}
	return profile.Image48, nil
}
