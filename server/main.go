package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"voice-server/config"
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
	Type string `json:"type"`
	Text string `json:"text"`
}

type WebSocketMessage struct {
	Messages []Message `json:"messages"`
}

type AudioMessage struct {
	Type      string    `json:"type"`
	AudioData string    `json:"audioData"`
	Messages  []Message `json:"messages"`
}

type TranscriptionResponse struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func callLLMAPI(messages []Message) (string, error) {
	fullMessages := []Message{
		{
			Role:    "developer",
			Content: config.SystemPrompt,
		},
	}
	fullMessages = append(fullMessages, messages...)

	chatRequest := ChatRequest{
		Model:    config.LLMModel,
		Messages: fullMessages,
	}

	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	log.Printf("Sending API request: %s", string(jsonData))

	resp, err := http.Post(config.LLMAPIUrl, "application/json", bytes.NewBuffer(jsonData))
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

func saveAudioBlob(audioData string) error {
	data, err := base64.StdEncoding.DecodeString(audioData)
	if err != nil {
		return fmt.Errorf("failed to decode base64 audio data: %v", err)
	}
	log.Printf("Processing audio blob of %d bytes", len(data))

	tempFile := filepath.Join(config.AudioDir, fmt.Sprintf("temp_%d.webm", time.Now().UnixNano()))
	log.Printf("Creating temporary file: %s", tempFile)

	err = os.WriteFile(tempFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary audio file: %v", err)
	}
	defer func() {
		if removeErr := os.Remove(tempFile); removeErr != nil {
			log.Printf("Warning: failed to remove temp file %s: %v", tempFile, removeErr)
		}
	}()

	outputFile := filepath.Join(config.AudioDir, config.QuestionAudioFile)
	log.Printf("Converting %s to %s using ffmpeg", tempFile, outputFile)

	cmd := exec.Command("ffmpeg", "-y", "-i", tempFile, "-ar", "22050", "-ac", "1", "-sample_fmt", "s16", outputFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("ffmpeg stdout: %s", stdout.String())
		log.Printf("ffmpeg stderr: %s", stderr.String())
		return fmt.Errorf("ffmpeg conversion failed: %v", err)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return fmt.Errorf("output file %s was not created", outputFile)
	}

	log.Printf("Audio blob successfully converted and saved to %s", outputFile)
	return nil
}

func transcribeAudio(audioPath string) (string, error) {
	log.Printf("Transcribing audio file: %s", audioPath)

	cmd := exec.Command(config.WhisperBinPath,
		"--no-prints",
		"--no-timestamps",
		"--language", config.WhisperLanguage,
		"--model", config.WhisperModel,
		"--file", audioPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("whisper-cli stdout: %s", stdout.String())
		log.Printf("whisper-cli stderr: %s", stderr.String())
		return "", fmt.Errorf("whisper-cli transcription failed: %v", err)
	}

	transcribedText := strings.TrimSpace(stdout.String())
	log.Printf("Transcribed text: %s", transcribedText)

	return transcribedText, nil
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
		var audioMsg AudioMessage
		if err := json.Unmarshal(message, &audioMsg); err == nil && audioMsg.Type == "audio" {
			log.Printf("Received audio message with %d bytes of audio data", len(audioMsg.AudioData))

			err = saveAudioBlob(audioMsg.AudioData)
			if err != nil {
				log.Printf("Error saving audio blob: %v", err)
				errorMsg := fmt.Sprintf("Error processing audio: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			log.Printf("Audio successfully processed and saved to question.wav")

			transcribedText, err := transcribeAudio(filepath.Join(config.AudioDir, config.QuestionAudioFile))
			if err != nil {
				log.Printf("Error transcribing audio: %v", err)
				errorMsg := fmt.Sprintf("Error transcribing audio: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			transcriptionResponse := TranscriptionResponse{
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

			userMessage := Message{
				Role:    "user",
				Content: transcribedText,
			}
			fullMessages := append(audioMsg.Messages, userMessage)

			llmOutput, err := callLLMAPI(fullMessages)
			if err != nil {
				log.Printf("LLM error: %v", err)
				errorMsg := fmt.Sprintf("Error processing with LLM: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			log.Printf("LLM: %s", llmOutput)

			piperCmd := fmt.Sprintf(`echo "%s" | %s/piper --model %s/%s --output_file %s/%s`,
				llmOutput, config.PiperDir, config.PiperDir, config.PiperModel, config.AudioDir, config.AnswerAudioFile)
			ttsCmd := exec.Command("bash", "-c", piperCmd)
			err = ttsCmd.Run()
			if err != nil {
				log.Printf("Piper error: %v", err)
				errorMsg := fmt.Sprintf("Error generating speech: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
				continue
			}

			assistantResponse := Response{
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
		var wsMessage WebSocketMessage
		err = json.Unmarshal(message, &wsMessage)
		if err != nil {
			log.Printf("Error parsing text message: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error parsing your request."))
			continue
		}

		llmOutput, err := callLLMAPI(wsMessage.Messages)
		if err != nil {
			log.Printf("LLM error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error processing your request."))
			continue
		}

		log.Printf("LLM: %s", llmOutput)

		piperCmd := fmt.Sprintf(`echo "%s" | %s/piper --model %s/%s --output_file %s/%s`,
			llmOutput, config.PiperDir, config.PiperDir, config.PiperModel, config.AudioDir, config.AnswerAudioFile)
		ttsCmd := exec.Command("bash", "-c", piperCmd)
		err = ttsCmd.Run()
		if err != nil {
			log.Printf("Piper error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error generating speech."))
			continue
		}

		response := Response{
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
