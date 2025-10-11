package models

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatChoice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
}

type Response struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type WebSocketMessage struct {
	Messages []Message `json:"messages"`
}

type AudioMessage struct {
	Type      string    `json:"type"`
	AudioData string    `json:"audioData"`
	Messages  []Message `json:"messages"`
}

type TranscriptionResponse struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

