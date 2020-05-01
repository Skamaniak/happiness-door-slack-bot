package handler

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
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

func writeResponse(response domain.SlackResponse, w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func HappinessDoorHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	response := domain.SlackResponse{Markdown: true, Text: "hello _Hello_ *HELLO!*"}
	err := writeResponse(response, w)
	if err != nil {
		log.Println("WARN: Failed to respond to request", err)
	}
}
