.PHONY: build clean test run help

BINARY_NAME=s3-directory-listing
BUILD_DIR=bin

.DEFAULT_GOAL := help

## build: Build executable
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/s3-directory-listing

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

## fmt: Format all Go files
fmt:
	gofmt -s -w .

## vet: Run go vet
vet:
	go vet ./...

## help: Display this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
