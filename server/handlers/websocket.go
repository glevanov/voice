package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

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

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Check if this is an audio message with chat history
		var audioMsg models.AudioMessage
		if err := json.Unmarshal(rawMessage, &audioMsg); err == nil && audioMsg.Type == "audio" {
			handleAudioMessage(conn, audioMsg)
			continue
		}

		log.Printf("Received: %s", rawMessage)

		// Parse the incoming message containing chat history
		var textMessage models.WebSocketMessage
		err = json.Unmarshal(rawMessage, &textMessage)
		if err != nil {
			log.Printf("Error parsing text message: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error parsing your request."))
			continue
		}

		handleTextMessage(conn, textMessage)
	}
}

func handleAudioMessage(conn *websocket.Conn, audioMsg models.AudioMessage) {
	log.Printf("Received audio message with %d bytes of audio data", len(audioMsg.AudioData))

	err := services.SaveAudioBlob(audioMsg.AudioData)
	if err != nil {
		log.Printf("Error saving audio blob: %v", err)
		errorMsg := fmt.Sprintf("Error processing audio: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		return
	}

	log.Printf("Audio successfully processed and saved to question.wav")

	transcription, err := services.TranscribeAudio(filepath.Join(config.AudioDir, config.QuestionAudioFile))
	if err != nil {
		log.Printf("Error transcribing audio: %v", err)
		errorMsg := fmt.Sprintf("Error transcribing audio: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		return
	}

	transcriptionResponse := models.TranscriptionResponse{
		Type: "user",
		Text: transcription,
	}
	transcriptionJSON, err := json.Marshal(transcriptionResponse)
	if err != nil {
		log.Printf("Error marshaling transcription response: %v", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, transcriptionJSON)
	if err != nil {
		log.Printf("Error sending transcription: %v", err)
		return
	}

	userMessage := models.Message{
		Role:    "user",
		Content: transcription,
	}
	conversationHistory := append(audioMsg.Messages, userMessage)

	assistantMessage, err := services.CallLLMAPI(conversationHistory)
	if err != nil {
		log.Printf("LLM error: %v", err)
		errorMsg := fmt.Sprintf("Error processing with LLM: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		return
	}

	log.Printf("LLM: %s", assistantMessage)

	err = services.GenerateSpeech(assistantMessage, config.AnswerAudioFile)
	if err != nil {
		log.Printf("Error generating speech: %v", err)
		errorMsg := fmt.Sprintf("Error generating speech: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		return
	}

	assistantResponse := models.Response{
		Type: "assistant",
		Text: assistantMessage,
	}
	assistantJSON, err := json.Marshal(assistantResponse)
	if err != nil {
		log.Printf("Error marshaling assistant response: %v", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, assistantJSON)
	if err != nil {
		log.Printf("Error sending assistant response: %v", err)
		return
	}
}

func handleTextMessage(conn *websocket.Conn, textMessage models.WebSocketMessage) {
	assistantMessage, err := services.CallLLMAPI(textMessage.Messages)
	if err != nil {
		log.Printf("LLM error: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error processing your request."))
		return
	}

	log.Printf("LLM: %s", assistantMessage)

	err = services.GenerateSpeech(assistantMessage, config.AnswerAudioFile)
	if err != nil {
		log.Printf("Error generating speech: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error generating speech."))
		return
	}

	response := models.Response{
		Type: "assistant",
		Text: assistantMessage,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}
