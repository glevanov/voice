package services

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"voice-server/config"
)

// TranscribeAudio transcribes an audio file using Whisper CLI
func TranscribeAudio(audioPath string) (string, error) {
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

	transcription := strings.TrimSpace(stdout.String())
	log.Printf("Transcribed text: %s", transcription)

	return transcription, nil
}
