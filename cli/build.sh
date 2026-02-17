#!/bin/bash
set -e

VERSION=${1:-"dev"}
OUTPUT_DIR="./dist"

echo "Building resume-cli version ${VERSION}..."

# Clean previous builds
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"

# Build for current platform
echo "Building for current platform..."
go build -ldflags="-s -w" -o "${OUTPUT_DIR}/resume-cli" .

echo "✓ Build complete: ${OUTPUT_DIR}/resume-cli"

# Optional: Build for multiple platforms
if [ "$2" = "multi" ]; then
    echo "Building for multiple platforms..."

    # macOS ARM64
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "${OUTPUT_DIR}/resume-cli-darwin-arm64" .

    # macOS AMD64
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "${OUTPUT_DIR}/resume-cli-darwin-amd64" .

    # Linux AMD64
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "${OUTPUT_DIR}/resume-cli-linux-amd64" .

    # Windows AMD64
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "${OUTPUT_DIR}/resume-cli-windows-amd64.exe" .

    echo "✓ Multi-platform builds complete"
fi

# Show file sizes
echo ""
echo "Build artifacts:"
ls -lh "${OUTPUT_DIR}"
