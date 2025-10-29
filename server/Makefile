.PHONY: run deps build build-linux fmt help

GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
BINARY_NAME=voice-server

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

# Build the binary for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v ./main.go

# Format code
fmt:
	go fmt ./...

help:
	@echo "Available targets:"
	@echo "  run          - Run the application"
	@echo "  deps         - Install dependencies"
	@echo "  build        - Build the binary"
	@echo "  build-linux  - Cross compile for Linux"
	@echo "  fmt          - Format source code"
	@echo "  help         - Show this help"
