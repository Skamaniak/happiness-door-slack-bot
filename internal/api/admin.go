package api

import (
	"github.com/Skamaniak/happiness-door-slack-bot/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AdminHandlers struct {
	service *service.SlackService
}

func NewAdminHandlers(service *service.SlackService) *AdminHandlers {
	return &AdminHandlers{service: service}
}

func (h *AdminHandlers) Ping(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Pong"))

	if err != nil {
		logrus.WithError(err).Warn("Failed to respond to ping request.")
	}
}

func (h *AdminHandlers) Health(w http.ResponseWriter, _ *http.Request) {
	err := h.service.HealthCheck()

	if err != nil {
		w.WriteHeader(503)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logrus.WithError(err).Warn("Failed to respond to health request.")
		}
	} else {
		w.WriteHeader(204)
	}
}
