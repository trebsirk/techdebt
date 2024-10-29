# Define project name and paths
PROJECT_NAME := techdebt
SRC_DIR := .
BUILD_DIR := bin
BIN := $(BUILD_DIR)/$(PROJECT_NAME)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go get github.com/go-git/go-git/v5

# Build the project
.PHONY: build
build: deps
	@echo "Building the project..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BIN) $(SRC_DIR)

# Run the project
.PHONY: run
run: build
	@echo "Running the project..."
	$(BIN)

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linting checks
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Help output
.PHONY: help
help:
	@echo "Makefile for the Git Commit Collector project"
	@echo ""
	@echo "Usage:"
	@echo "  make deps       - Install dependencies"
	@echo "  make build      - Build the project"
	@echo "  make run        - Build and run the project"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make fmt        - Format the code"
	@echo "  make lint       - Run linting checks"
	@echo "  make help       - Show this help message"
