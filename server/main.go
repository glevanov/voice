package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	LLM_MODEL   = "gemma3:4b"
	PIPER_MODEL = "sv_SE-nst-medium.onnx"
	PORT        = ":3002"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

type Response struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received: %s", message)

		// Run Ollama LLM
		cmd := exec.Command("ollama", "run", LLM_MODEL, string(message))
		stdout, err := cmd.Output()
		if err != nil {
			log.Printf("LLM error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error processing your request."))
			continue
		}

		llmOutput := strings.TrimSpace(string(stdout))
		log.Printf("LLM: %s", llmOutput)

		// Run Piper TTS
		piperCmd := fmt.Sprintf(`echo "%s" | ../piper/piper --model ../piper/%s --output_file ../audio/answer.wav`, llmOutput, PIPER_MODEL)
		ttsCmd := exec.Command("bash", "-c", piperCmd)
		err = ttsCmd.Run()
		if err != nil {
			log.Printf("Piper error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error generating speech."))
			continue
		}

		// Send response
		response := Response{
			Text:  llmOutput,
			Audio: "answer.wav",
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, responseJSON)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
	}
}

func serveAudio(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
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

	// Construct the full path to the audio file
	audioPath := filepath.Join("../audio", filename)

	// Check if file exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		http.Error(w, "Audio file not found", http.StatusNotFound)
		return
	}

	// Set appropriate content type for WAV files
	w.Header().Set("Content-Type", "audio/wav")

	// Serve the file
	http.ServeFile(w, r, audioPath)
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", serveAudio)

	log.Printf("WebSocket server starting on port %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
