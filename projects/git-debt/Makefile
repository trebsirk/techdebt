# Define directories
PROJECT_NAME := git-debt
SRC_DIR := .
BIN_DIR := bin
BIN := $(BIN_DIR)/$(PROJECT_NAME)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go get github.com/go-git/go-git/v5

# Build the project
.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(SRC_DIR)/main.go

# Run the project
.PHONY: run
run: build
	$(BIN)

# Clean up build artifacts
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
