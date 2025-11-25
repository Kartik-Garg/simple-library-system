# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

# Install build dependencies (git for Go modules, ca-certificates for HTTPS)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first (better layer caching - only re-download if these change)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary (CGO_ENABLED=0 for portability, -ldflags strips debug info for smaller size)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o main ./cmd/api/main.go

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies (ca-certificates for HTTPS, tzdata for timezone support)
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security (UID/GID 1000)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy only the binary from build stage (keeps image small)
COPY --from=builder /app/main .

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Run as non-root user (security best practice)
USER appuser

# Document which port the app uses
EXPOSE 8080

# Health check - Docker will ping /health endpoint every 30s
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]

