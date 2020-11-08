package client

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"time"
)

type SlackClient struct {
	api           *slack.Client
	userIconCache *cache.Cache
}

func NewSlackClient() *SlackClient {
	token := viper.GetString(conf.SlackToken)
	api := slack.New(token)
	userIconCache := cache.New(6*time.Hour, time.Hour)

	return &SlackClient{
		api:           api,
		userIconCache: userIconCache,
	}
}

func (c *SlackClient) GetUserIconUrl(userId string) (string, error) {
	if val, found := c.userIconCache.Get(userId); found {
		logrus.WithField("userId", userId).Debug("User profile icon taken from cache")
		return val.(string), nil
	}

	profile, err := c.api.GetUserProfile(userId, false)
	if err != nil {
		return "", err
	}
	iconUrl := profile.Image48
	c.userIconCache.Set(userId, iconUrl, cache.DefaultExpiration)
	logrus.WithField("userId", userId).Debug("User profile icon saved to cache")
	return iconUrl, nil
}

func (c *SlackClient) PostMessage(channelID string, msg slack.Blocks) (string, error) {
	_, messageTS, err := c.api.PostMessage(channelID, slack.MsgOptionBlocks(msg.BlockSet...))
	return messageTS, err
}

func (c *SlackClient) ReplaceMessage(channelID string, messageTS string, msg slack.Blocks) error {
	_, _, _, err := c.api.UpdateMessage(channelID, messageTS, slack.MsgOptionBlocks(msg.BlockSet...))
	return err
}

func (c *SlackClient) CanPostMessage(channelID string) (bool, error) {
	cnl, err := c.api.GetConversationInfo(channelID, false)
	if err != nil {
		if err.Error() == "channel_not_found" {
			// Not a member and channel is private
			return false, nil
		}
		return false, err
	}

	return !cnl.IsPrivate || cnl.IsMember, nil
}
