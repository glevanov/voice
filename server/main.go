package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	LLM_MODEL   = "google/gemma-3n-e4b"
	LLM_API_URL = "http://localhost:1234/v1/chat/completions"
	PIPER_MODEL = "sv_SE-nst-medium.onnx"
	PORT        = ":3002"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatChoice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
}

type Response struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type WebSocketMessage struct {
	Messages []Message `json:"messages"`
}

func callLLMAPI(messages []Message) (string, error) {
	// Prepend the developer message
	fullMessages := []Message{
		{
			Role:    "developer",
			Content: "You are a helpful and friendly conversation partner. Always answer in Swedish, as if you were talking to a friend. Answer as if you are speaking, avoiding using emojis, special characters, formatting or comments in your responses. Focus on natural language and a personal tone.",
		},
	}
	fullMessages = append(fullMessages, messages...)

	chatRequest := ChatRequest{
		Model:    LLM_MODEL,
		Messages: fullMessages,
	}

	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	log.Printf("Sending API request: %s", string(jsonData))

	resp, err := http.Post(LLM_API_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Received API response: %s", string(body))

	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResponse.Choices[0].Message.Content, nil
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

		// Parse the incoming message containing chat history
		var wsMessage WebSocketMessage
		err = json.Unmarshal(message, &wsMessage)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error parsing your request."))
			continue
		}

		// Call LLM API with full chat history
		llmOutput, err := callLLMAPI(wsMessage.Messages)
		if err != nil {
			log.Printf("LLM error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error processing your request."))
			continue
		}

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
