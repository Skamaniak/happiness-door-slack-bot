package handler

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/slack-go/slack"
	"log"
	"net/http"
)

func writeResponse(response slack.Msg, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("WARN: failed to marshal slack message to JSON")
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println("WARN: failed to write slack message to response")
	}
	return err
}

func HappinessDoorHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("WARN: Failed to parse request", err)
	} else {
		meetingName := r.Form.Get("text")
		message := domain.CreateInitMessage(meetingName)
		err := writeResponse(message, w)
		if err != nil {
			log.Println("WARN: Failed to respond to request", err)
		}
	}
}
