#!/bin/bash
# Start all device-anomaly-detection-demo services using docker-compose
set -e

echo "Starting device-anomaly-detection-demo services..."

# Check if config.yaml exists for me-restful-mapper
if [ ! -f "./me-restful-mapper/config.yaml" ]; then
  echo "Warning: me-restful-mapper/config.yaml not found. Please ensure it exists before starting."
fi

# Check if /etc/kubeedge directory exists
if [ ! -d "/etc/kubeedge" ]; then
  echo "Warning: /etc/kubeedge directory not found. me-restful-mapper may not work properly."
fi

# Start services
docker-compose up -d

echo "All services started successfully!"
echo "  - room-simulator: http://localhost:8000"
echo "  - external-anomaly-detector: http://localhost:6804"
echo "  - me-restful-mapper: running"
echo ""
echo "Use 'docker-compose logs -f' to view logs"
echo "Use './stop.sh' to stop all services"
