package client

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/patrickmn/go-cache"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"log"
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
		log.Println(fmt.Sprintf("INFO: User profile icon for user %s taken from cache", userId))
		return val.(string), nil
	}

	profile, err := c.api.GetUserProfile(userId, false)
	if err != nil {
		return "", err
	}
	iconUrl := profile.Image48
	c.userIconCache.Set(userId, iconUrl, cache.DefaultExpiration)
	log.Println(fmt.Sprintf("INFO: User profile icon for user %s saved to cache", userId))
	return iconUrl, nil
}
