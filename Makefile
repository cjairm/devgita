.PHONY: all build-darwin-arm64 build-darwin-amd64 build-linux-amd64 clean test

# Build all platform binaries
all: build-darwin-arm64 build-darwin-amd64 build-linux-amd64

# macOS ARM64 (Apple Silicon - M1, M2, M3 chips)
build-darwin-arm64:
	@echo "Building for macOS ARM64 (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -o devgita-darwin-arm64 .
	@echo "✓ devgita-darwin-arm64 built successfully"

# macOS AMD64 (Intel chips)
build-darwin-amd64:
	@echo "Building for macOS AMD64 (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -o devgita-darwin-amd64 .
	@echo "✓ devgita-darwin-amd64 built successfully"

# Linux AMD64 (Debian/Ubuntu)
build-linux-amd64:
	@echo "Building for Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build -o devgita-linux-amd64 .
	@echo "✓ devgita-linux-amd64 built successfully"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f devgita-darwin-arm64 devgita-darwin-amd64 devgita-linux-amd64
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "✓ Tests passed"

# Run code quality checks
lint:
	@echo "Running code quality checks..."
	go vet ./...
	go fmt ./...
	@echo "✓ Code quality checks passed"

# Build for current platform only
build:
	@echo "Building for current platform..."
	go build -o devgita .
	@echo "✓ devgita built successfully"

# Help
help:
	@echo "Available targets:"
	@echo "  all                - Build all platform binaries"
	@echo "  build-darwin-arm64 - Build for macOS ARM64 (Apple Silicon)"
	@echo "  build-darwin-amd64 - Build for macOS AMD64 (Intel)"
	@echo "  build-linux-amd64  - Build for Linux AMD64"
	@echo "  build              - Build for current platform"
	@echo "  clean              - Remove build artifacts"
	@echo "  test               - Run test suite"
	@echo "  lint               - Run code quality checks"
	@echo "  help               - Show this help message"
