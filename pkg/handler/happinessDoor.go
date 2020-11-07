package handler

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Handlers struct {
	service *service.SlackService
}

func NewHandlers(service *service.SlackService) *Handlers {
	return &Handlers{service: service}
}

func logRequest(r *http.Request) {
	if requestBytes, err := httputil.DumpRequest(r, true); err != nil {
		logrus.WithError(err).Warn("Failed to parse request")
	} else {
		logrus.Debug(string(requestBytes))
	}
}

func (h *Handlers) Initiation(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	defer func() { _ = r.Body.Close() }()

	slash, err := slack.SlashCommandParse(r)
	if err != nil {
		logrus.WithError(err).Warn("Failed to parse request")
		return
	}

	resp, err := h.service.InitiateHappinessDoor(slash.Text, slash.ChannelID)
	if err != nil {
		logrus.WithError(err).Warn("Failed to create new happiness door")
		return
	}
	if resp != nil {
		err = writeResponse(resp, w)
		if err != nil {
			logrus.WithError(err).Warn("Failed to respond to request")
		}
	}
}

func (h *Handlers) Vote(_ http.ResponseWriter, r *http.Request) {
	logRequest(r)
	defer func() { _ = r.Body.Close() }()

	err := r.ParseForm()
	if err != nil {
		logrus.WithError(err).Warn("Failed to parse form")
		return
	}

	logrus.WithField("Form", r.Form).Debug("Form parsed")
	payload, _ := url.QueryUnescape(r.Form.Get("payload"))
	var result domain.InteractiveResponse

	err = json.Unmarshal([]byte(payload), &result)
	if err != nil {
		logrus.WithError(err).Warn("Failed to parse response from payload parameter")
		return
	}

	err = h.service.IncrementVoting(result)
	if err != nil {
		logrus.WithError(err).Warn("Failed to increment voting")
	}
}

func toJson(v interface{}) []byte {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		logrus.WithError(err).Warn("Failed to marshal slack message to JSON")
	}
	return jsonBytes
}

func writeResponse(response *slack.Msg, w http.ResponseWriter) error {
	w.Header().
		Set("Content-Type", "application/json")

	jsonBytes := toJson(response)
	_, err := w.Write(jsonBytes)
	return err
}
