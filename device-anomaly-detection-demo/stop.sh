#!/bin/bash
# Stop all device-anomaly-detection-demo services using docker-compose
set -e

echo "Stopping device-anomaly-detection-demo services..."

docker-compose down

echo "All services stopped successfully!"
