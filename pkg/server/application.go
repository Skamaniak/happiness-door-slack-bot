package server

import (
	"fmt"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, r *http.Request) {
	message := "Hello World!"
	_, err := fmt.Fprintln(w, message)
	if err != nil {
		log.Println("WARN: Failed to respond to request", err)
	}
}
