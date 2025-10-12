package main

import (
	"log"
	"net/http"

	"voice-server/config"
	"voice-server/handlers"
)

func main() {
	http.HandleFunc("/ws", handlers.HandleWebSocket)
	http.HandleFunc("/", handlers.ServeAudio)

	log.Printf("WebSocket server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
