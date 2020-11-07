package client

import (
	"bytes"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"net/http"
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

func (_ *SlackClient) SendToSlack(url string, request []byte) {
	logrus.WithField("Url", url).Debug("Sending response to Slack")
	r, err := http.Post(url, "application/json", bytes.NewBuffer(request))
	defer func() { _ = r.Body.Close() }()
	if err != nil {
		logrus.WithError(err).Warn("Failed to send http request to response URL")
	}
	if r.StatusCode >= 400 {
		logrus.WithFields(logrus.Fields{"Response": r}).Warn("Got unexpected response from Slack")
	}
}

func (c *SlackClient) PostMessage(channelID string, msg slack.Blocks) (string, error) {
	_, messageTS, err := c.api.PostMessage(channelID, slack.MsgOptionBlocks(msg.BlockSet...))
	return messageTS, err
}

func (c *SlackClient) PostEphemeralMessage(channelID, userID string, msg slack.Blocks) error {
	_, err := c.api.PostEphemeral(channelID, userID, slack.MsgOptionBlocks(msg.BlockSet...))
	return err
}

func (c *SlackClient) ReplaceMessage(channelID string, messageTS string, msg slack.Blocks) error {
	_, _, _, err := c.api.UpdateMessage(channelID, messageTS, slack.MsgOptionBlocks(msg.BlockSet...))
	return err
}

func (c *SlackClient) IsBotMember(channelID string) (bool, error) {
	cnl, err := c.api.GetConversationInfo(channelID, false)
	if err != nil {
		return false, err
	}
	return cnl.IsMember, nil
}
