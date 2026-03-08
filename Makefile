.PHONY: build clean test run help

BINARY_NAME=s3-directory-listing
BUILD_DIR=bin

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  test     - Run tests"
	@echo "  run      - Build and run the binary"
	@echo "  help     - Show this help message"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/s3-directory-listing

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

run: build
	@$(BUILD_DIR)/$(BINARY_NAME)
