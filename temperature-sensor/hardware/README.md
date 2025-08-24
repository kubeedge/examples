# Modbus Simulator

## Description
This is a modbus simulator built using gomodbus dependencies, currently only supports TCP connections

## Quick Start

### 1. Build the image
```
docker build -t temperature-sensor:v1.0 .
```

### 2. Deploy the simulator
```
kubectl apply -f deploy.yaml
```