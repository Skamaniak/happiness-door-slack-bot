package server

import (
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"log"
	"net/http"
	"net/http/httputil"
)

func HappinessDoorHandler(w http.ResponseWriter, r *http.Request) {
	if requestBytes, err := httputil.DumpRequest(r, true); err != nil {
		log.Println("Failed to parse request", err)
	} else {
		log.Println(string(requestBytes))
	}

	response := domain.SlackResponse{Markdown: true, Text: "hello _Hello_ *HELLO!*"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("WARN: Failed to respond to request", err)
	}
}
