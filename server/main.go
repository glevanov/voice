package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
	"voice-server/config"
	"voice-server/models"
	"voice-server/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
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

		// Check if this is an audio message with chat history
		var audioMsg models.AudioMessage
		if err := json.Unmarshal(message, &audioMsg); err == nil && audioMsg.Type == "audio" {
			log.Printf("Received audio message with %d bytes of audio data", len(audioMsg.AudioData))

			err = services.SaveAudioBlob(audioMsg.AudioData)
			if err != nil {
				log.Printf("Error saving audio blob: %v", err)
				errorMsg := fmt.Sprintf("Error processing audio: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			log.Printf("Audio successfully processed and saved to question.wav")

			transcribedText, err := services.TranscribeAudio(filepath.Join(config.AudioDir, config.QuestionAudioFile))
			if err != nil {
				log.Printf("Error transcribing audio: %v", err)
				errorMsg := fmt.Sprintf("Error transcribing audio: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			transcriptionResponse := models.TranscriptionResponse{
				Type: "user",
				Text: transcribedText,
			}
			responseJSON, err := json.Marshal(transcriptionResponse)
			if err != nil {
				log.Printf("Error marshaling transcription response: %v", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, responseJSON)
			if err != nil {
				log.Printf("Error sending transcription: %v", err)
				break
			}

			userMessage := models.Message{
				Role:    "user",
				Content: transcribedText,
			}
			fullMessages := append(audioMsg.Messages, userMessage)

			llmOutput, err := services.CallLLMAPI(fullMessages)
			if err != nil {
				log.Printf("LLM error: %v", err)
				errorMsg := fmt.Sprintf("Error processing with LLM: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			log.Printf("LLM: %s", llmOutput)

			err = services.GenerateSpeech(llmOutput, config.AnswerAudioFile)
			if err != nil {
				log.Printf("Error generating speech: %v", err)
				errorMsg := fmt.Sprintf("Error generating speech: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			assistantResponse := models.Response{
				Type: "assistant",
				Text: llmOutput,
			}
			assistantResponseJSON, err := json.Marshal(assistantResponse)
			if err != nil {
				log.Printf("Error marshaling assistant response: %v", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, assistantResponseJSON)
			if err != nil {
				log.Printf("Error sending assistant response: %v", err)
				break
			}

			continue
		}

		log.Printf("Received: %s", message)

		// Parse the incoming message containing chat history
		var wsMessage models.WebSocketMessage
		err = json.Unmarshal(message, &wsMessage)
		if err != nil {
			log.Printf("Error parsing text message: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error parsing your request."))
			continue
		}

		llmOutput, err := services.CallLLMAPI(wsMessage.Messages)
		if err != nil {
			log.Printf("LLM error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error processing your request."))
			continue
		}

		log.Printf("LLM: %s", llmOutput)

		err = services.GenerateSpeech(llmOutput, config.AnswerAudioFile)
		if err != nil {
			log.Printf("Error generating speech: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error generating speech."))
			continue
		}

		response := models.Response{
			Type: "assistant",
			Text: llmOutput,
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
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", serveAudio)

	log.Printf("WebSocket server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
