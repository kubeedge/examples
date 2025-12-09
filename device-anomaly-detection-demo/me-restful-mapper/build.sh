#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Define variables
BUILD_DIR="build"
OUTPUT_BINARY="main"
DOCKERFILE="Dockerfile_nostream"
VERSION_FILE=".version"
IMAGE_NAME="huajuan6848/me-restful-mapper:$(cat $VERSION_FILE)"

# Step 1: Read version
if [ ! -f "$VERSION_FILE" ]; then
    echo "Version file not found!"
    exit 1
fi
VERSION=$(cat "$VERSION_FILE")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Step 2: Clean previous build artifacts
if [ -d "$BUILD_DIR" ]; then
    echo "Cleaning previous build artifacts..."
    rm -rf "$BUILD_DIR"
fi
mkdir "$BUILD_DIR"

# Step 3: Build the Go binary
echo "Building the Go binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X 'main.version=$VERSION' -X 'main.buildDate=$BUILD_DATE'" -o "$BUILD_DIR/$OUTPUT_BINARY" cmd/main.go

# Step 4: Build the Docker image
echo "Building the Docker image..."
docker build -f "$DOCKERFILE" -t "$IMAGE_NAME" .

# Step 5: Done
echo "Build process completed successfully. Docker image: $IMAGE_NAME"