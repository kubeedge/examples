#!/bin/bash
# Stop the room-simulator Docker container
set -e

CONTAINER_NAME="room-simulator"

# Check if container exists
if [ $(docker ps -aq -f name=$CONTAINER_NAME) ]; then
    echo "Stopping and removing container $CONTAINER_NAME ..."
    docker rm -f $CONTAINER_NAME
    echo "Container $CONTAINER_NAME stopped successfully."
else
    echo "Container $CONTAINER_NAME is not running."
fi
