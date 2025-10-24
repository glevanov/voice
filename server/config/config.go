package config

const (
	// LLM Configuration
	LLMModel  = "openai/gpt-oss-20b"
	LLMAPIUrl = "http://localhost:1234/v1/chat/completions"
	PiperModel = "sv_SE-nst-medium.onnx"
	// "ggml-base.bin", "ggml-medium.bin", "ggml-large-v3.bin"
	WhisperModel = "ggml-large-v3.bin"

	// Server Configuration
	Port = ":3002"

	// Path Configuration
	AudioDir       = "../audio"
	PiperDir       = "../bin/piper"
	ModelsDir      = "../models"

	// Audio Processing Configuration
	QuestionAudioFile = "question.wav"
	AnswerAudioFile   = "answer.wav"

	// API URL
	WhisperAPIUrl = "http://localhost:4002"
	PiperAPIUrl = "http://localhost:4001"

	// System Prompt
	SystemPrompt = `You are a helpful and friendly conversation partner.
	Your purpose is to help the user practice their Swedish speaking skills.
	Always answer in Swedish, even if user asks you to speak another language.
	Respond only in plain text,	avoiding using emojis, special characters, formatting or comments in your responses.
	Focus on natural language and a personal tone.
	If you notice unnatural phrasing or severe grammatical mistakes, also offer corrections in your response.`
)
