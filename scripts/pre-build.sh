#!/bin/bash

# Pre-build script for GoReleaser
# This script runs before building the binaries

set -e

echo "🔨 Starting pre-build process..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# Check if we have a tag
if [ -z "$GORELEASER_CURRENT_TAG" ]; then
    echo "⚠️  Warning: No git tag found, using development version"
else
    echo "🏷️  Building version: $GORELEASER_CURRENT_TAG"
fi

# Validate Go version
GO_VERSION=$(go version | awk '{print $3}')
echo "⚡ Go version: $GO_VERSION"

# Check if required Go version is met (1.21+)
REQUIRED_VERSION="go1.21"
if ! go version | grep -q "go1.2[1-9]\|go.[3-9][0-9]\|go[2-9]"; then
    echo "❌ Error: Go version 1.21 or higher is required"
    echo "Current version: $GO_VERSION"
    exit 1
fi

# Run go mod tidy
echo "🧹 Running go mod tidy..."
go mod tidy

# Run tests
echo "🧪 Running tests..."
go test -v ./...

# Check for any linting issues
echo "🔍 Running linting checks..."
if command -v golangci-lint >/dev/null 2>&1; then
    golangci-lint run
else
    echo "⚠️  golangci-lint not found, skipping linting"
fi

# Create build info
echo "📝 Creating build information..."
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

echo "Build Time: $BUILD_TIME"
echo "Git Commit: $GIT_COMMIT"
echo "Git Branch: $GIT_BRANCH"

echo "✅ Pre-build process completed successfully!" 