package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func HappinessDoorHandler(w http.ResponseWriter, r *http.Request) {
	var requestStr string

	if requestBytes, err := httputil.DumpRequest(r, true); err != nil {
		requestStr = "Failed to parse request"
		log.Println(requestStr, err)
	} else {
		requestStr = string(requestBytes)
		log.Println(requestStr)
	}

	_, err := fmt.Fprintf(w, "Request body: %v", requestStr)
	if err != nil {
		log.Println("WARN: Failed to respond to request", err)
	}
}
