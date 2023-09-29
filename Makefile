# Copyright (c) Tetrate, Inc 2022 All Rights Reserved.

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: lint test build release clean

# Variables
BIN_DIR     		?= ./bin
BINARY_NAME 		:= lictl
GIT_REPO        := boeboe/lictl
RELEASE_VERSION ?= v0.1.0

LINTER := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2


lint: ## Run linter on all source code
	@echo "Running linter..."
	@go run $(LINTER) run --verbose ./...


test: ## Run all tests recursively
	@echo "Running tests..."
	@go test -v ./...


build: lint test ## Build the project
	@echo "Building the binary..."
	@go build -o $(BIN_DIR)/$(BINARY_NAME) .
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-x86_64 .
	GOOS=linux GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-macos-arm64 .



release: build ## Create a GitHub release and upload the binary
	@which gh >/dev/null || (echo "gh is not installed" && exit 1)
	@echo "Checking if release $(RELEASE_VERSION) already exists..."
	@if gh release view $(RELEASE_VERSION) -R $(GIT_REPO) > /dev/null 2>&1; then \
		echo "Release $(RELEASE_VERSION) exists. Deleting it..."; \
		gh release delete $(RELEASE_VERSION) -R $(GIT_REPO) --yes; \
	fi
	@echo "Creating a new release on GitHub..."
	gh release create $(RELEASE_VERSION) \
		$(BIN_DIR)/$(BINARY_NAME)-x86_64 \
		$(BIN_DIR)/$(BINARY_NAME)-arm64 \
		$(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe \
		$(BIN_DIR)/$(BINARY_NAME)-macos-amd64 \
		$(BIN_DIR)/$(BINARY_NAME)-macos-arm64 \
		--title "Release $(RELEASE_VERSION)" --notes "Release notes for $(RELEASE_VERSION)" --repo $(GIT_REPO)


# Clean up build output artifacts
clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf $(BIN_DIR)
