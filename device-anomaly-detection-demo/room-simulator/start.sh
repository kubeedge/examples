#!/bin/bash
# Start the room-simulator Docker container
set -e

IMAGE_NAME="huajuan6848/room-simulator:1.0.0"
CONTAINER_NAME="room-simulator"
PORT=8000

# Check and create bridge network if needed
NETWORK="device-anomaly-detection-demo"
if ! docker network ls --format '{{.Name}}' | grep -wq "$NETWORK"; then
    echo "Docker network $NETWORK does not exist. Creating..."
    docker network create --driver bridge $NETWORK
fi

# Remove existing container if exists
if [ $(docker ps -aq -f name=$CONTAINER_NAME) ]; then
    echo "Removing existing container $CONTAINER_NAME ..."
    docker rm -f $CONTAINER_NAME
fi

echo "Starting container $CONTAINER_NAME ..."
docker run -d --name $CONTAINER_NAME --network $NETWORK -p $PORT:8000 $IMAGE_NAME

echo "room-simulator is running at http://localhost:$PORT"
