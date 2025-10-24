package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"voice-server/config"
)

func TranscribeAudio(audioPath string) (string, error) {
	log.Printf("Transcribing audio file: %s", audioPath)

	whisperURL := fmt.Sprintf("%s/transcribe", config.WhisperAPIUrl)

	resp, err := http.Get(whisperURL)
	if err != nil {
		return "", fmt.Errorf("Failed to call Whisper service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Whisper service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read Whisper response: %v", err)
	}

	transcription := strings.TrimSpace(string(body))
	log.Printf("Transcribed text: %s", transcription)

	return transcription, nil
}
