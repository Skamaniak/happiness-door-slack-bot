package server

import (
	"fmt"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, r *http.Request) {
	message := "Hello World!"
	_, err := fmt.Print(w, message)
	if err != nil {
		log.Fatal("Failed to respond to request", err)
	}
}
