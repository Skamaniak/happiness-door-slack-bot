package handler

import (
	"bytes"
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/slack-go/slack"
	"log"
	"net/http"
	"net/http/httputil"
)

func toJson(v interface{}) []byte {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Println("WARN: failed to marshal slack message to JSON")
	}
	return jsonBytes
}

func logRequest(r *http.Request) {
	if requestBytes, err := httputil.DumpRequest(r, true); err != nil {
		log.Println("Failed to parse request", err)
	} else {
		log.Println(string(requestBytes))
	}
}

func writeResponse(response slack.Msg, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes := toJson(response)
	_, err := w.Write(jsonBytes)
	return err
}

func HappinessDoorHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	logRequest(r)

	if s, err := slack.SlashCommandParse(r); err != nil {
		log.Println("WARN: Failed to parse request", err)
	} else {
		message := domain.CreateInitMessage(s.Text)
		err := writeResponse(message, w)
		if err != nil {
			log.Println("WARN: Failed to respond to request", err)
		}
	}
}

func Interaction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	logRequest(r)

	if s, err := slack.SlashCommandParse(r); err != nil {
		log.Println("WARN: Failed to parse request", err)
	} else {
		responseUrl := s.ResponseURL
		log.Println("Got response URL", responseUrl)

		jsonBytes := toJson(domain.CreateResultMessage())

		_, err := http.Post(responseUrl, "application/json", bytes.NewBuffer(jsonBytes))
		if err != nil {
			log.Println("WARN: Failed to send http request to response URL")
		}
	}
}
