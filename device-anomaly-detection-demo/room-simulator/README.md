# Room Simulator

A Python application that simulates a room with light devices and a brightness sensor, providing RESTful API interfaces for device control and status querying.

## Features

- **Two Light Devices**: `light1` and `light2`, each can be on (1) or off (0)
- **Brightness Sensor**: Weighted sum of the actual state of each light
  - `brightness = light1.actual_state × 1 + light2.actual_state × 2`
  - 0 = dark (both lights off)
  - 1 = dim (light1 on, light2 off)
  - 2 = moderate (light1 off, light2 on)
  - 3 = bright (both lights on)
- **Fault Injection**: Inject faults into any light device, causing its actual working state to be opposite to its set state
- **RESTful API**: Complete REST interface for device control and status querying
- **Logging**: Detailed logging of all device state changes

## Requirements

- Python >= 3.10
- Dependencies:
  - FastAPI >= 0.104.0
  - uvicorn[standard] >= 0.24.0
  - pydantic >= 2.5.0

## Installation

1. Clone or download the project to your local machine

2. Install dependencies:
```bash
pip install -e .
```

Or using uv (recommended):
```bash
uv pip install -e .
```

## Usage

Start the room simulator service:

```bash
python main.py
```

The service will run on `http://0.0.0.0:8000`.

You can also specify a port:
```bash
python main.py --port 9000
```

Or specify both host and port:
```bash
python main.py --host 127.0.0.1 --port 9000
```

### Command Line Options

- `--port`: Port to run the server on (default: 8000)
- `--host`: Host to run the server on (default: 0.0.0.0)

## API Documentation

After starting the service, you can access interactive API documentation at:

- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## API Endpoints

### 1. Get Welcome Message

```
GET /
```

Returns basic API information.

### 2. Get All Device Status

```
GET /status
```

Returns complete status information for all devices in the room.

**Response Example**:
```json
{
  "light1": {
    "light": "light1",
    "set_state": 0,
    "actual_state": 0,
    "has_fault": false
  },
  "light2": {
    "light": "light2",
    "set_state": 0,
    "actual_state": 0,
    "has_fault": false
  },
  "brightness_sensor": {
    "brightness": 0,
    "description": "dark",
    "light1_actual_state": 0,
    "light2_actual_state": 0
  }
}
```

### 3. Get Light State

```
GET /lights/{light_name}
```

**Path Parameters**:
- `light_name`: Light name (`light1` or `light2`)

**Response**: Plain text representing the set state
```
1
```
(Returns `0` for off or `1` for on)

### 4. Set Light State

```
PUT /lights/{light_name}
```

**Path Parameters**:
- `light_name`: Light name (`light1` or `light2`)

**Request Body**: Plain text
```
1
```
- `0` = off, `1` = on

**Response**: Plain text
```
ok
```

### 5. Inject Fault

```
POST /lights/{light_name}/fault
```

**Path Parameters**:
- `light_name`: Light name (`light1` or `light2`)

**Request Body**:
```
0
```
- `0` = Remove fault
- `1` = Inject fault

**Response**: JSON
```json
{
    "light": "light1",
    "fault_injected": true,
    "set_state": 0,
    "actual_state": 1,
    "brightness": 0
}
```

### 6. Get Brightness Sensor Value

```
GET /brightness
```

**Response**: Plain text representing brightness level
```
1
```
(Returns `0` for dark, `1` for dim, `2` for moderate, or `3` for bright)

## Usage Examples

### Using curl

```bash
# Get all status
curl http://localhost:8000/status

# Get light1 state
curl http://localhost:8000/lights/light1

# Turn on light1
curl -X PUT http://localhost:8000/lights/light1 \
  -H "Content-Type: text/plain" \
  -d '1'

# Turn on light2
curl -X PUT http://localhost:8000/lights/light2 \
  -H "Content-Type: text/plain" \
  -d '1'

# Check brightness sensor
curl http://localhost:8000/brightness

# Inject fault into light1
curl -X POST http://localhost:8000/lights/light1/fault \
  -H "Content-Type: text/plain" \
  -d '1'

# Remove fault from light1
curl -X POST http://localhost:8000/lights/light1/fault \
  -H "Content-Type: text/plain" \
  -d '0'

# Turn off light2
curl -X PUT http://localhost:8000/lights/light2 \
  -H "Content-Type: text/plain" \
  -d '0'
```

### Using Python requests

```python
import requests

base_url = "http://localhost:8000"

# Get all status
response = requests.get(f"{base_url}/status")
print(response.json())

# Get light1 state
response = requests.get(f"{base_url}/lights/light1")
print(f"Light1 state: {response.text}")  # Returns "0" or "1"

# Turn on light1
response = requests.put(
    f"{base_url}/lights/light1",
    data="1",
    headers={"Content-Type": "text/plain"}
)
print(response.text)  # Returns "ok"

# Inject fault
response = requests.post(
    f"{base_url}/lights/light1/fault",
    data="1",
    headers={"Content-Type": "text/plain"}
)
print(response.json())

# Check brightness
response = requests.get(f"{base_url}/brightness")
print(f"Brightness: {response.text}")  # Returns "0", "1", "2", or "3"
```

## How It Works

### Light Device

- Each light device has a **set state** and an **actual state**
- Under normal conditions: set state = actual state
- When a fault is injected: actual state = 1 - set state (state inversion)

### Brightness Sensor

- Brightness value is based on the **actual state** of lights in the room (not the set state)
- Calculation formula: `brightness = light1.actual_state × 1 + light2.actual_state × 2`
- This weighted design makes each combination of light states produce a unique brightness value, enabling the anomaly detector to unambiguously identify which device is faulty
- Therefore, even if a faulty light is set to on but actually off, the brightness sensor correctly reflects the actual brightness

### Fault Injection Example

Scenario: light1 is set to on (1), then a fault is injected

1. Set state: 1 (on)
2. After fault injection
3. Actual state: 0 (off, due to fault causing state inversion)
4. Brightness sensor calculates brightness based on actual state (0)

## Logging

All device state changes are logged to:
- Console output
- `room_simulator.log` file

Log example:
```
2025-12-01 10:00:00 - __main__ - INFO - Light device light1 initialized, state: off
2025-12-01 10:00:05 - __main__ - INFO - Light device light1 state changed: 0 -> 1 (on)
2025-12-01 10:00:05 - __main__ - INFO - Brightness sensor brightness_sensor state changed: 0(dark) -> 1(dim)
2025-12-01 10:00:10 - __main__ - WARNING - Light device light1 fault injected
2025-12-01 10:00:10 - __main__ - INFO - Light device light1 actual working state: off
2025-12-01 10:00:10 - __main__ - INFO - Brightness sensor brightness_sensor state changed: 1(dim) -> 0(dark)
```

## Project Structure

```
room-simulator/
├── main.py              # Main program file with all classes and API endpoints
├── pyproject.toml       # Project configuration and dependencies
├── README.md            # Project documentation
└── room_simulator.log   # Log file (generated at runtime)
```

## Architecture

### Class Design

1. **Light Class**: Represents a single light device
   - Manages set state and actual state
   - Supports fault injection
   - Logs state changes

2. **BrightnessSensor Class**: Represents the brightness sensor
   - Calculates brightness as a weighted sum: `light1.actual_state × 1 + light2.actual_state × 2`
   - Logs brightness changes

3. **Room Class**: Room manager
   - Contains two Light instances and one BrightnessSensor instance
   - Coordinates interactions between devices
   - Provides unified interface for API calls

### Fault Model

Fault injection simulates device hardware faults, causing the device's actual behavior to differ from expected:
- When set to on, a faulty device is actually off
- When set to off, a faulty device is actually on
- This fault affects the brightness sensor readings

## License

MIT License

## Authors

Room Simulator Development Team
