package services

import (
	"fmt"
	"os/exec"
	"voice-server/config"
)

// GenerateSpeech takes text and generates speech using Piper TTS,
// saving the output to the specified output file in the audio directory.
func GenerateSpeech(text string, outputFilename string) error {
	piperCmd := fmt.Sprintf(`echo "%s" | %s/piper --model %s/%s --output_file %s/%s`,
		text, config.PiperDir, config.PiperDir, config.PiperModel, config.AudioDir, outputFilename)

	ttsCmd := exec.Command("bash", "-c", piperCmd)
	err := ttsCmd.Run()
	if err != nil {
		return fmt.Errorf("piper TTS error: %w", err)
	}

	return nil
}

