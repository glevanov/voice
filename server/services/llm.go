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

func CallLLMAPI(messages []models.Message) (string, error) {
	conversationMessages := []models.Message{
		{
			Role:    "developer",
			Content: config.SystemPrompt,
		},
	}
	conversationMessages = append(conversationMessages, messages...)

	chatRequest := models.ChatRequest{
		Model:    config.LLMModel,
		Messages: conversationMessages,
	}

	requestJSON, err := json.Marshal(chatRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	log.Printf("Sending API request: %s", string(requestJSON))

	httpResp, err := http.Post(config.LLMAPIUrl, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("error making API request: %v", err)
	}
	defer httpResp.Body.Close()

	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Received API response: %s", string(responseBody))

	var chatResponse models.ChatResponse
	err = json.Unmarshal(responseBody, &chatResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResponse.Choices[0].Message.Content, nil
}
