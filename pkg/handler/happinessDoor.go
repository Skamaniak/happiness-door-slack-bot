package handler

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/slack-go/slack"
	"log"
	"net/http"
	"net/http/httputil"
)

func logRequest(r *http.Request) {
	if requestBytes, err := httputil.DumpRequest(r, true); err != nil {
		log.Println("Failed to parse request", err)
	} else {
		log.Println(string(requestBytes))
	}
}

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
	logRequest(r)

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

func Interaction(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	//message := domain.CreateResultMessage()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"replace_original\": \"true\",\"text\": \"Thanks for your request, we'll process it and get back to you.\"}"))
	//err := writeResponse(message, w)
	//if err != nil {
	//	log.Println("WARN: Failed to respond to request", err)
	//}
}
