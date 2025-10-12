package config

const (
	// LLM Configuration
	LLMModel  = "google/gemma-3n-e4b"
	LLMAPIUrl = "http://localhost:1234/v1/chat/completions"

	// TTS Configuration
	PiperModel = "sv_SE-nst-medium.onnx"

	// Server Configuration
	Port = ":3002"

	// Path Configuration
	AudioDir       = "../audio"
	PiperDir       = "../piper"
	WhisperBinPath = "../../whisper.cpp/build/bin/whisper-cli"
	WhisperModel   = "../../whisper.cpp/models/ggml-base.bin"

	// Audio Processing Configuration
	QuestionAudioFile = "question.wav"
	AnswerAudioFile   = "answer.wav"

	// Whisper Configuration
	WhisperLanguage = "auto"

	// System Prompt
	SystemPrompt = "You are a helpful and friendly conversation partner. Always answer in Swedish, as if you were talking to a friend. Answer as if you are speaking, avoiding using emojis, special characters, formatting or comments in your responses. Focus on natural language and a personal tone."
)
