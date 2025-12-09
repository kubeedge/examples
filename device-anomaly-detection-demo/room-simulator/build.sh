#!/bin/bash
# Build the Docker image for room-simulator
set -e

IMAGE_NAME="huajuan6848/room-simulator:1.0.0"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "Building Docker image: $IMAGE_NAME ..."
docker build -t $IMAGE_NAME .
echo "Build complete."
