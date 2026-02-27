#!/bin/bash
# Start the external-anomaly-detector Docker container
set -e

IMAGE_NAME="huajuan6848/external-anomaly-detector:1.0.0"
CONTAINER_NAME="external-anomaly-detector"
PORT=6804

# Check and create bridge network if needed
NETWORK="device-anomaly-detection-demo"
if ! docker network ls --format '{{.Name}}' | grep -wq "$NETWORK"; then
    echo "Docker network $NETWORK does not exist. Creating..."
    docker network create --driver bridge $NETWORK
fi

# Prepare kubeconfig mount if exists
KUBECONFIG_HOST="$HOME/.kube/config"
KUBECONFIG_CONTAINER="/root/.kube/config"
KUBECONFIG_MOUNT=""
if [ -f "$KUBECONFIG_HOST" ]; then
    echo "Mounting $KUBECONFIG_HOST to container $KUBECONFIG_CONTAINER"
    KUBECONFIG_MOUNT="-v $KUBECONFIG_HOST:$KUBECONFIG_CONTAINER:ro"
fi

# Remove existing container if exists
if [ $(docker ps -aq -f name=$CONTAINER_NAME) ]; then
    echo "Removing existing container $CONTAINER_NAME ..."
    docker rm -f $CONTAINER_NAME
fi

echo "Starting container $CONTAINER_NAME ..."
docker run -d --name $CONTAINER_NAME --network $NETWORK -p $PORT:6804 $KUBECONFIG_MOUNT $IMAGE_NAME

echo "external-anomaly-detector is running at http://localhost:$PORT"