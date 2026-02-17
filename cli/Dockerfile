FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o resume-cli .

# Runtime image
FROM alpine:latest

# Install texlive for PDF generation
RUN apk add --no-cache texlive-full

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/resume-cli /usr/local/bin/resume-cli

# Default command
ENTRYPOINT ["resume-cli"]
CMD ["serve", "--port", "8080"]
