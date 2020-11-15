package client

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"time"
)

type SlackUser struct {
	UID     string
	Email   string
	Name    string
	IconUrl string
}

type SlackClient struct {
	api              *slack.Client
	iconByIdCache    *cache.Cache
	iconByEmailCache *cache.Cache
}

func NewSlackClient() *SlackClient {
	token := viper.GetString(conf.SlackToken)
	api := slack.New(token)
	iconByIdCache := cache.New(6*time.Hour, time.Hour)
	iconByEmailCache := cache.New(6*time.Hour, time.Hour)

	return &SlackClient{
		api:              api,
		iconByIdCache:    iconByIdCache,
		iconByEmailCache: iconByEmailCache,
	}
}

func (c *SlackClient) GetUserByEmail(e string) (SlackUser, error) {
	if val, found := c.iconByEmailCache.Get(e); found {
		logrus.WithField("userEmail", e).Debug("User taken from cache")
		return val.(SlackUser), nil
	}

	user, err := c.api.GetUserByEmail(e)
	if err != nil {
		return SlackUser{}, err
	}
	slackUser := userFromProfile(user.ID, user.Profile)
	c.iconByIdCache.Set(e, slackUser, cache.DefaultExpiration)
	logrus.WithField("userEmail", e).Debug("User saved to cache")
	return slackUser, nil
}

func (c *SlackClient) GetUserById(uID string) (SlackUser, error) {
	if val, found := c.iconByIdCache.Get(uID); found {
		logrus.WithField("userId", uID).Debug("User taken from cache")
		return val.(SlackUser), nil
	}

	profile, err := c.api.GetUserProfile(uID, false)
	if err != nil {
		return SlackUser{}, err
	}
	slackUser := userFromProfile(uID, *profile)
	c.iconByIdCache.Set(uID, slackUser, cache.DefaultExpiration)
	logrus.WithField("userId", uID).Debug("User saved to cache")
	return slackUser, nil
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

func userFromProfile(uID string, p slack.UserProfile) SlackUser {
	return SlackUser{UID: uID, Email: p.Email, Name: p.DisplayName, IconUrl: p.Image48}
}
