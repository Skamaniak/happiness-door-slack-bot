package handler

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"log"
	"net/http"
)

func writeResponse(response string, w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func HappinessDoorHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("WARN: Failed to parse request", err)
	} else {
		meetingName := r.Form.Get("text")
		log.Println("Creating happiness door for ", meetingName)

		response := domain.CreateInitMessage(meetingName)
		err := writeResponse(response, w)
		if err != nil {
			log.Println("WARN: Failed to respond to request", err)
		}
	}
}
