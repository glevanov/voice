.PHONY: build run clean deps build-linux fmt help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=voice-server
BINARY_UNIX=$(BINARY_NAME)_unix

# Run the application
run:
	$(GOCMD) run main.go

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./main.go

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f output.wav

# Cross compilation for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./main.go

# Format code
fmt:
	go fmt ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  run          - Run the application"
	@echo "  deps         - Install dependencies"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  build-linux  - Cross compile for Linux"
	@echo "  fmt          - Format source code"
	@echo "  help         - Show this help"
