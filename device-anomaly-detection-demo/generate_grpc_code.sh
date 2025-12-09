#!/bin/bash

# Define paths for the projects and proto file
PROTO_FILE="./api/external_anomaly_detector.proto"
GO_PROJECT="./me-restful-mapper"
PYTHON_PROJECT="./external-anomaly-detector"

# Check if the proto file exists
if [ ! -f "$PROTO_FILE" ]; then
  echo "Proto file not found: $PROTO_FILE"
  exit 1
fi

# Generate gRPC code for the Golang project
echo "Generating gRPC code for Golang project..."
# Ensure output directory exists
mkdir -p "./$GO_PROJECT/api/externalanomalydetectorpb"
protoc --go_out="./$GO_PROJECT/api/externalanomalydetectorpb" \
       --go_opt=paths=source_relative \
       --go_opt=Mexternal_anomaly_detector.proto=github.com/kubeedge/me-restful-mapper/api/externalanomalydetectorpb \
       --go-grpc_out="./$GO_PROJECT/api/externalanomalydetectorpb" \
       --go-grpc_opt=paths=source_relative \
       --go-grpc_opt=Mexternal_anomaly_detector.proto=github.com/kubeedge/me-restful-mapper/api/externalanomalydetectorpb \
       --proto_path="./api" "$PROTO_FILE"
if [ $? -ne 0 ]; then
  echo "Failed to generate gRPC code for Golang project."
  exit 1
fi

echo "Golang gRPC code generated successfully."

# Generate gRPC code for the Python project
echo "Generating gRPC code for Python project..."
# Ensure output directory exists
mkdir -p "$PYTHON_PROJECT/externalanomalydetectorpb"
python -m grpc_tools.protoc -I./api --python_out="$PYTHON_PROJECT/externalanomalydetectorpb" --grpc_python_out="$PYTHON_PROJECT/externalanomalydetectorpb" "$PROTO_FILE"
if [ $? -ne 0 ]; then
  echo "Failed to generate gRPC code for Python project."
  exit 1
fi

echo "Python gRPC code generated successfully."

echo "gRPC code generation completed for both projects."