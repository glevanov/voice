package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"voice-server/config"
	"voice-server/models"
)

// CallLLMAPI sends messages to the LLM API and returns the response
func CallLLMAPI(messages []models.Message) (string, error) {
	fullMessages := []models.Message{
		{
			Role:    "developer",
			Content: config.SystemPrompt,
		},
	}
	fullMessages = append(fullMessages, messages...)

	chatRequest := models.ChatRequest{
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

	var chatResponse models.ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResponse.Choices[0].Message.Content, nil
}

