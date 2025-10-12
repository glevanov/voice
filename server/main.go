package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"voice-server/config"
	"voice-server/handlers"
)

func serveAudio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	// Extract filename from URL path
	filename := strings.TrimPrefix(r.URL.Path, "/")
	if filename == "" {
		http.Error(w, "No filename provided", http.StatusBadRequest)
		return
	}

	audioPath := filepath.Join(config.AudioDir, filename)

	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		http.Error(w, "Audio file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "audio/wav")

	http.ServeFile(w, r, audioPath)
}

func main() {
	http.HandleFunc("/ws", handlers.HandleWebSocket)
	http.HandleFunc("/", serveAudio)

	log.Printf("WebSocket server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
