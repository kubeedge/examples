#!/bin/bash
# Build Docker image for external-anomaly-detector
set -e
IMAGE_NAME="huajuan6848/external-anomaly-detector:1.0.0"
docker build -t $IMAGE_NAME .
echo "Build complete: $IMAGE_NAME"