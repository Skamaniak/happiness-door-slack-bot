package cmd

import (
	"../pkg/server"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", server.HelloServer)
	err := http.ListenAndServe(":8080", nil)

	log.Fatal("Failed to create HTTP server", err)
}
