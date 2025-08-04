# Multi-stage build for Kavach CLI
# Stage 1: Build the binary
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
    -X github.com/Gkemhcs/kavach-cli/internal/version.Version=${VERSION:-dev} \
    -X github.com/Gkemhcs/kavach-cli/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
    -X github.com/Gkemhcs/kavach-cli/internal/version.GitCommit=$(git rev-parse HEAD 2>/dev/null || echo unknown) \
    -X github.com/Gkemhcs/kavach-cli/internal/version.GitBranch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)" \
    -a -installsuffix cgo -o kavach ./main.go

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S kavach && \
    adduser -u 1001 -S kavach -G kavach

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/kavach .

# Create config directory
RUN mkdir -p /app/config && \
    chown -R kavach:kavach /app

# Switch to non-root user
USER kavach

# Set environment variables
ENV KAVACH_CONFIG=/app/config

# Expose port (if needed for API)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./kavach"]

# Default command
CMD ["--help"] 