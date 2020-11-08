package api

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/ws"
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

// TODO make the parsing nicer, also tie the user identity and hdID with the socket initiation - that is when the token is verified.
// TODO Otherwise anyone an vote as anyone and for any happiness door if they passed auth for one particular HD.

// HelloFromClient is a method that handles messages from the app client.
func (h *WSHandlers) vote(_ *ws.Socket, data interface{}) {
	byteData, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to convert data to bytes.")
	}
	ia := domain.InteractiveResponse{}
	err = json.Unmarshal(byteData, &ia)
	if err != nil {
		logrus.WithError(err).Error("Failed to convert bytes to interactive action dto.")
	}

	err = h.service.IncrementVoting(ia)
	if err != nil {
		logrus.WithError(err).Error("Failed to increment voting.")
	}
}
