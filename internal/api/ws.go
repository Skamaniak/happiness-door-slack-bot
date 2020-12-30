package api

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/domain"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/service"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/ws"
	"github.com/sirupsen/logrus"
)

type WSHandlers struct {
	service *service.SlackService
}

func NewWSHandlers(service *service.SlackService) *WSHandlers {
	return &WSHandlers{service: service}
}

func (h *WSHandlers) CreateVoteHandler() (ws.Event, ws.Handler) {
	return "VoteAction", h.vote
}

// TODO make the parsing nicer
// HelloFromClient is a method that handles messages from the app client.
func (h *WSHandlers) vote(s *ws.Socket, data interface{}) {
	byteData, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to convert data to bytes.")
	}
	a := domain.Action{}
	err = json.Unmarshal(byteData, &a)
	if err != nil {
		logrus.WithError(err).Error("Failed to convert bytes to interactive action dto.")
	}

	err = h.service.WebVoting(s.AuthContext.UserEmail, a)
	if err != nil {
		logrus.WithError(err).Error("Failed to increment voting.")
	}
}
