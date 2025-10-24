package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"voice-server/config"
)

type PiperRequest struct {
	Text string `json:"text"`
}

func GenerateSpeech(text string) error {
	requestBody := PiperRequest{
		Text: text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		config.PiperAPIUrl,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("piper HTTP request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("piper server returned status: %d", resp.StatusCode)
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	outputPath := filepath.Join(config.AudioDir, config.AnswerAudioFile)
	err = os.WriteFile(outputPath, audioData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write audio file: %w", err)
	}

	return nil
}
