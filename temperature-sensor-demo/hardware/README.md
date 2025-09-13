# Modbus Simulator

## Description
This is a modbus simulator built using gomodbus dependencies, currently only supports TCP connections

## Quick Start

### 1. Build the image
```
docker build -t temperature-sensor:v1.0 .
```
### 2. Upload image to the edge node

### 3. Deploy the simulator to the edge
```
kubectl apply -f ./resource/deploy.yaml
```