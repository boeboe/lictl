# Variables
BINARY_NAME=licli
BIN_DIR=./bin

# Build the project
build:
	@echo "Building the binary..."
	@go build -o $(BIN_DIR)/$(BINARY_NAME) .

# Run linter on all source code
lint:
	@echo "Running linter..."
	@golint ./...

# Run all tests recursively
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean up build output artifacts
clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf $(BIN_DIR)

.PHONY: build lint test clean
