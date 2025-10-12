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

func SaveAudioBlob(audioData string) error {
	decodedAudio, err := base64.StdEncoding.DecodeString(audioData)
	if err != nil {
		return fmt.Errorf("failed to decode base64 audio data: %v", err)
	}
	log.Printf("Processing audio blob of %d bytes", len(decodedAudio))

	tempFilePath := filepath.Join(config.AudioDir, fmt.Sprintf("temp_%d.webm", time.Now().UnixNano()))
	log.Printf("Creating temporary file: %s", tempFilePath)

	err = os.WriteFile(tempFilePath, decodedAudio, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary audio file: %v", err)
	}
	defer func() {
		if removeErr := os.Remove(tempFilePath); removeErr != nil {
			log.Printf("Warning: failed to remove temp file %s: %v", tempFilePath, removeErr)
		}
	}()

	outputFilePath := filepath.Join(config.AudioDir, config.QuestionAudioFile)
	log.Printf("Converting %s to %s using ffmpeg", tempFilePath, outputFilePath)

	cmd := exec.Command("ffmpeg", "-y", "-i", tempFilePath, "-ar", "22050", "-ac", "1", "-sample_fmt", "s16", outputFilePath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("ffmpeg stdout: %s", stdout.String())
		log.Printf("ffmpeg stderr: %s", stderr.String())
		return fmt.Errorf("ffmpeg conversion failed: %v", err)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		return fmt.Errorf("output file %s was not created", outputFilePath)
	}

	log.Printf("Audio blob successfully converted and saved to %s", outputFilePath)
	return nil
}
