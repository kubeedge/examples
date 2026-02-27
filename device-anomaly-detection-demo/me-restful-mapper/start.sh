#!/bin/bash
# Start me-restful-mapper service using Docker

set -e

IMAGE="huajuan6848/me-restful-mapper:1.0.0"
CONTAINER_NAME="me-restful-mapper"
KUBEEDGE_HOST_PATH="/etc/kubeedge"
KUBEEDGE_CONTAINER_PATH="/etc/kubeedge"
CONFIG_HOST_PATH="./config.yaml"
CONFIG_CONTAINER_PATH="/kubeedge/config.yaml"
MAIN_BIN="/kubeedge/main"

# Check if config.yaml exists
if [ ! -f "$CONFIG_HOST_PATH" ]; then
  echo "Config file $CONFIG_HOST_PATH not found. Please make sure config.yaml is generated."
  exit 1
fi

# Check if /etc/kubeedge directory exists
if [ ! -d "$KUBEEDGE_HOST_PATH" ]; then
  echo "Host path $KUBEEDGE_HOST_PATH not found. Please make sure it is mounted."
  exit 1
fi

# Check and create bridge network if needed
NETWORK="device-anomaly-detection-demo"
if ! docker network ls --format '{{.Name}}' | grep -wq "$NETWORK"; then
    echo "Docker network $NETWORK does not exist. Creating..."
    docker network create --driver bridge $NETWORK
fi

# Start the container
exec docker run --rm \
  -d \
  --name $CONTAINER_NAME \
  --network $NETWORK \
  -v "$KUBEEDGE_HOST_PATH:$KUBEEDGE_CONTAINER_PATH" \
  -v "$CONFIG_HOST_PATH:$CONFIG_CONTAINER_PATH" \
  $IMAGE $MAIN_BIN --config-file $CONFIG_CONTAINER_PATH
