# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates for module downloads
RUN apk add --no-cache git ca-certificates tzdata

# Set build environment
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy go mod and sum files for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN go build -ldflags="-w -s" -a -installsuffix cgo -o main ./main.go

# Runtime stage - using distroless for minimal size and security
FROM gcr.io/distroless/static-debian12:nonroot

# Copy timezone data from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder stage
COPY --from=builder /build/main /app/main

# Copy ca-certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Run as non-root user (distroless nonroot user)
USER nonroot:nonroot

# Run the application
ENTRYPOINT ["/app/main"]
