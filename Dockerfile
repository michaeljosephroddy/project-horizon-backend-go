# Build the Go binary
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN go build -o app .

# Create minimal image for running the binary
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port (if needed)
EXPOSE 9095 

# Run the binary
CMD ["./app"]
