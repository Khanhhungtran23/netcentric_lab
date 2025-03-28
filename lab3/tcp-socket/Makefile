.PHONY: all build clean run-server run-client test

# Default build flags
LDFLAGS = -s -w
BUILD_DIR = build
SERVER_DIR = cmd/server
CLIENT_DIR = cmd/client
BIN_DIR = bin

# Get Go version
GO_VERSION := $(shell go version | awk '{print $$3}')

# Default target
all: clean build

# Build everything
build: build-server build-client

# Build server
build-server:
	@echo "Building server with $(GO_VERSION)..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/server $(SERVER_DIR)
	@echo "Server built successfully!"

# Build client
build-client:
	@echo "Building client with $(GO_VERSION)..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/client $(CLIENT_DIR)
	@echo "Client built successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@mkdir -p $(BIN_DIR)
	@mkdir -p data
	@mkdir -p files
	@mkdir -p downloads
	@echo "Clean complete!"

# Run server
run-server: build-server
	@echo "Starting server..."
	@$(BIN_DIR)/server

# Run client
run-client: build-client
	@echo "Starting client..."
	@$(BIN_DIR)/client

# Create some sample text files for testing file download
create-sample-files:
	@echo "Creating sample files..."
	@mkdir -p files
	@echo "This is a sample text file for testing file download functionality." > files/sample.txt
	@echo "This is another sample text file with a different content." > files/another.txt
	@echo "Lorem ipsum dolor sit amet, consectetur adipiscing elit." > files/lorem.txt
	@echo "Sample files created!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests completed!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Formatting complete!"

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build everything (default)"
	@echo "  build        - Build server and client"
	@echo "  build-server - Build only the server"
	@echo "  build-client - Build only the client"
	@echo "  clean        - Clean build artifacts"
	@echo "  run-server   - Build and run the server"
	@echo "  run-client   - Build and run the client"
	@echo "  create-sample-files - Create sample text files for testing"
	@echo "  test         - Run tests"
	@echo "  fmt          - Format code"
	@echo "  help         - Show this help message"