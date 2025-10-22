.PHONY: build-piper build-whisper build-images compose-up compose-down

build-piper:
	docker build -t voice-whisper ./docker/voice-piper

build-whisper:
	docker build -t voice-whisper ./docker/voice-whisper

build-images: build-piper build-whisper

compose-up:
	docker-compose -f docker/docker-compose.yml up -d

compose-down:
	docker-compose -f docker/docker-compose.yml down