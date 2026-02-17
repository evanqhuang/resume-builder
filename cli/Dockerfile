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

# Install texlive for PDF generation and fontconfig for font discovery
RUN apk add --no-cache texlive-full fontconfig \
    && mkdir -p /etc/fonts/conf.d \
    && echo '<?xml version="1.0"?><!DOCTYPE fontconfig SYSTEM "urn:fontconfig:fonts.dtd"><fontconfig><dir>/usr/share/texmf-dist/fonts/opentype</dir></fontconfig>' > /etc/fonts/conf.d/99-texlive.conf \
    && fc-cache -f

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/resume-cli /usr/local/bin/resume-cli

# Default command
ENTRYPOINT ["resume-cli"]
CMD ["serve", "--port", "8080"]
