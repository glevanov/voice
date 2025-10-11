package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"voice-server/config"
)

// SaveAudioBlob decodes base64 audio data, saves it as a temporary file,
// and converts it to WAV format using ffmpeg
func SaveAudioBlob(audioData string) error {
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

