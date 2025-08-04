# Makefile for Kavach CLI
# This file provides common development and build tasks

.PHONY: help build test clean install lint format release snapshot docker-build docker-push

# Default target
help:
	@echo "🔖 Kavach CLI - Available Commands"
	@echo ""
	@echo "📦 Build Commands:"
	@echo "  build        - Build the CLI for current platform"
	@echo "  build-all    - Build for all supported platforms"
	@echo "  install      - Install the CLI locally"
	@echo ""
	@echo "🧪 Test Commands:"
	@echo "  test         - Run all tests"
	@echo "  test-race    - Run tests with race detection"
	@echo "  test-coverage - Run tests with coverage"
	@echo ""
	@echo "🔧 Development Commands:"
	@echo "  lint         - Run linter"
	@echo "  format       - Format code"
	@echo "  tidy         - Tidy go modules"
	@echo ""
	@echo "🚀 Release Commands:"
	@echo "  release      - Build and release (requires tag)"
	@echo "  snapshot     - Build snapshot release"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-push  - Push Docker image"
	@echo ""
	@echo "🧹 Utility Commands:"
	@echo "  clean        - Clean build artifacts"
	@echo "  version      - Show current version"

# Build variables
BINARY_NAME=kavach
VERSION=$(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X github.com/Gkemhcs/kavach-cli/internal/version.Version=${VERSION} -X github.com/Gkemhcs/kavach-cli/internal/version.BuildTime=${BUILD_TIME} -X github.com/Gkemhcs/kavach-cli/internal/version.GitCommit=${GIT_COMMIT} -X github.com/Gkemhcs/kavach-cli/internal/version.GitBranch=${GIT_BRANCH}"

# Build for current platform
build:
	@echo "🔨 Building ${BINARY_NAME} for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@go build ${LDFLAGS} -o ${BINARY_NAME} ./main.go
	@echo "✅ Build complete: ./${BINARY_NAME}"

# Build for all supported platforms
build-all:
	@echo "🔨 Building for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 ./main.go
	@GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 ./main.go
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 ./main.go
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 ./main.go
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe ./main.go
	@echo "✅ All builds complete in dist/ directory"

# Install locally
install:
	@echo "📦 Installing ${BINARY_NAME}..."
	@go install ${LDFLAGS} ./main.go
	@echo "✅ Installed: $(shell go env GOPATH)/bin/${BINARY_NAME}"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with race detection
test-race:
	@echo "🧪 Running tests with race detection..."
	@go test -race -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

# Run linter
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Format code
format:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@goimports -w .

# Tidy go modules
tidy:
	@echo "🧹 Tidying go modules..."
	@go mod tidy
	@go mod verify

# Build and release (requires git tag)
release:
	@echo "🚀 Building release..."
	@goreleaser release --clean

# Build snapshot release
snapshot:
	@echo "📸 Building snapshot..."
	@goreleaser build --snapshot --clean

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t kavach-cli:${VERSION} .
	@docker tag kavach-cli:${VERSION} kavach-cli:latest

# Push Docker image
docker-push:
	@echo "🐳 Pushing Docker image..."
	@docker push kavach-cli:${VERSION}
	@docker push kavach-cli:latest

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf dist/
	@rm -f ${BINARY_NAME}
	@rm -f ${BINARY_NAME}.exe
	@rm -f coverage.out
	@rm -f coverage.html
	@go clean -cache
	@echo "✅ Clean complete"

# Show current version
version:
	@echo "🔖 Current version: ${VERSION}"
	@echo "📅 Build time: ${BUILD_TIME}"
	@echo "🔗 Git commit: ${GIT_COMMIT}"
	@echo "🌿 Git branch: ${GIT_BRANCH}"

# Development setup
setup:
	@echo "🔧 Setting up development environment..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/goreleaser/goreleaser@latest
	@echo "✅ Development environment setup complete"

# Check if all dependencies are available
check-deps:
	@echo "🔍 Checking dependencies..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed"; exit 1; }
	@command -v git >/dev/null 2>&1 || { echo "❌ Git is not installed"; exit 1; }
	@command -v goreleaser >/dev/null 2>&1 || { echo "⚠️  GoReleaser not found, run 'make setup'"; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "⚠️  golangci-lint not found, run 'make setup'"; }
	@echo "✅ All required dependencies are available" 