# Modbus Mapper Plugin

This project is a Modbus protocol temperature sensor management plugin based on KubeEdge, which implements data collection, status reporting and control functions for temperature sensors.

## Functional Features

- Support Modbus TCP and RTU communication modes
- Implement temperature data collection and reporting
- Support temperature alarm threshold setting
- Provide RESTful API interface for device control and data query
- Support real-time monitoring of device status


## Prerequisites

- KubeEdge 1.21.0+
- Go 1.22.9+
- Modbus device

## Quick Start

### 1. Configure mapper plugin
- Configure /crds/temperature-instance.yaml
```
1. replace the ip with your modbus device ip 
2. replace the nodeName with your edge node name
```
- Configure /resource/deployment.yaml
```
1. replace the nodeName with your edge node name
```
### 2. Build mapper plugin image at Edge
```
make docker-build
```
### 3. Deploy mapper and crds
```
make deploy
```
### 4. Verify plugin function
#### 4.1 View the synchronization of the reported field in the twins field:
```
kubectl get device -o yaml
```
#### 4.2 Access the REST API of the mapper plugin on the edge(default port:7777)/cloud node(default port:30077)
- Health check
```
curl <ip>:<port>/api/v1/ping
```

- Get device temperature
```
curl <ip>:<port>/api/v1/device/default/temperature-instance/temperature
```

- Set device temperature
```
curl <ip>:<port>/api/v1/devicemethod/default/temperature-instance/UpdateTemperature/temperature/{data} &&echo
```
- Control the switch of device
```
curl <ip>:<port>/api/v1/devicemethod/default/temperature-instance/SwitchControl/temperature-switch/{data} &&echo
```

#### 4.3 View mapper logs

### 5. Uninstall the plugin
```
make clean
```