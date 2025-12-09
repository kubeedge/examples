"""
Room Simulator
--------------
Simulates a room with two light devices and a brightness sensor.
Supports controlling lights and querying sensor status via RESTful API.
"""

import logging
import argparse
from typing import Dict
from fastapi import FastAPI, HTTPException, Response, Request
from pydantic import BaseModel
import uvicorn

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('room_simulator.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

class Light:
    """
    Light Device Class
    
    Attributes:
        name: Device name
        state: Device state (0=off, 1=on)
        has_fault: Whether the device has a fault
        _actual_state: Actual working state
    """
    
    def __init__(self, name: str):
        """Initialize light device"""
        self.name = name
        self._state = 0  # Set state
        self.has_fault = False  # Whether fault exists
        logger.info(f"Light device {self.name} initialized, state: off")
    
    @property
    def state(self) -> int:
        """Get the set state"""
        return self._state
    
    @property
    def actual_state(self) -> int:
        """Get actual working state (considering faults)"""
        if self.has_fault:
            # When faulty, actual state is opposite of set state
            return 1 - self._state
        return self._state
    
    def set_state(self, state: int) -> None:
        """
        Set light state
        
        Args:
            state: Target state (0=off, 1=on)
        """
        if state not in [0, 1]:
            raise ValueError("State must be 0 (off) or 1 (on)")
        
        old_state = self._state
        self._state = state
        
        state_text = "on" if state == 1 else "off"
        logger.info(f"Light device {self.name} state changed: {old_state} -> {state} ({state_text})")
        
        if self.has_fault:
            actual_text = "off" if self.actual_state == 0 else "on"
            logger.warning(f"Light device {self.name} has fault, actual state: {actual_text}")
    
    def inject_fault(self, fault: bool) -> None:
        """
        Inject or remove fault
        
        Args:
            fault: True=inject fault, False=remove fault
        """
        self.has_fault = fault
        fault_status = "injected" if fault else "removed"
        logger.warning(f"Light device {self.name} fault {fault_status}")
        
        if fault:
            actual_text = "off" if self.actual_state == 0 else "on"
            logger.info(f"Light device {self.name} actual working state: {actual_text}")


class BrightnessSensor:
    """
    Brightness Sensor Class
    
    Automatically calculates brightness based on the number of lights in the room:
    - 0 = dark (no lights)
    - 1 = moderate (one light)
    - 2 = bright (two lights)
    """
    
    def __init__(self, name: str = "brightness_sensor"):
        """Initialize brightness sensor"""
        self.name = name
        self._brightness = 0
        logger.info(f"Brightness sensor {self.name} initialized, current brightness: dark")
    
    def update_brightness(self, active_lights_count: int) -> None:
        """
        Update brightness based on number of active lights
        
        Args:
            active_lights_count: Number of lights turned on
        """
        old_brightness = self._brightness
        self._brightness = active_lights_count
        
        brightness_map = {0: "dark", 1: "moderate", 2: "bright"}
        if old_brightness != self._brightness:
            logger.info(
                f"Brightness sensor {self.name} state changed: "
                f"{old_brightness}({brightness_map[old_brightness]}) -> "
                f"{self._brightness}({brightness_map[self._brightness]})"
            )
    
    @property
    def brightness(self) -> int:
        """Get current brightness value"""
        return self._brightness


class Room:
    """
    Room Class
    
    Manages all devices in the room, including two lights and one brightness sensor.
    """
    
    def __init__(self):
        """Initialize room and its devices"""
        self.light1 = Light("light1")
        self.light2 = Light("light2")
        self.brightness_sensor = BrightnessSensor("brightness_sensor")
        logger.info("Room initialized with devices: light1, light2, brightness_sensor")
        self._update_brightness()
    
    def _update_brightness(self) -> None:
        """Update brightness sensor based on actual state of lights"""
        active_count = self.light1.actual_state + self.light2.actual_state
        self.brightness_sensor.update_brightness(active_count)
    
    def set_light(self, light_name: str, state: int) -> Dict:
        """
        Set the state of a specified light
        
        Args:
            light_name: Light name ("light1" or "light2")
            state: Target state (0=off, 1=on)
            
        Returns:
            Dictionary containing operation result
        """
        if light_name == "light1":
            light = self.light1
        elif light_name == "light2":
            light = self.light2
        else:
            raise ValueError(f"Unknown light device: {light_name}")
        
        light.set_state(state)
        self._update_brightness()
        
        return {
            "light": light_name,
            "set_state": state,
            "actual_state": light.actual_state,
            "has_fault": light.has_fault,
            "brightness": self.brightness_sensor.brightness
        }
    
    def get_light_status(self, light_name: str) -> Dict:
        """
        Get the status of a specified light
        
        Args:
            light_name: Light name ("light1" or "light2")
            
        Returns:
            Dictionary containing light status
        """
        if light_name == "light1":
            light = self.light1
        elif light_name == "light2":
            light = self.light2
        else:
            raise ValueError(f"Unknown light device: {light_name}")
        
        return {
            "light": light_name,
            "set_state": light.state,
            "actual_state": light.actual_state,
            "has_fault": light.has_fault
        }
    
    def inject_fault(self, light_name: str, fault: bool) -> Dict:
        """
        Inject or remove fault from a specified light
        
        Args:
            light_name: Light name ("light1" or "light2")
            fault: True=inject fault, False=remove fault
            
        Returns:
            Dictionary containing operation result
        """
        if light_name == "light1":
            light = self.light1
        elif light_name == "light2":
            light = self.light2
        else:
            raise ValueError(f"Unknown light device: {light_name}")
        
        light.inject_fault(fault)
        self._update_brightness()
        
        return {
            "light": light_name,
            "fault_injected": fault,
            "set_state": light.state,
            "actual_state": light.actual_state,
            "brightness": self.brightness_sensor.brightness
        }
    
    def get_brightness(self) -> Dict:
        """
        Get brightness sensor status
        
        Returns:
            Dictionary containing brightness information
        """
        brightness_map = {0: "dark", 1: "moderate", 2: "bright"}
        return {
            "brightness": self.brightness_sensor.brightness,
            "description": brightness_map[self.brightness_sensor.brightness],
            "light1_actual_state": self.light1.actual_state,
            "light2_actual_state": self.light2.actual_state
        }
    
    def get_all_status(self) -> Dict:
        """
        Get status of all devices in the room
        
        Returns:
            Dictionary containing all device statuses
        """
        return {
            "light1": self.get_light_status("light1"),
            "light2": self.get_light_status("light2"),
            "brightness_sensor": self.get_brightness()
        }


# Pydantic models for API request validation
class LightStateRequest(BaseModel):
    """Light state setting request"""
    state: int
    
    class Config:
        json_schema_extra = {
            "example": {
                "state": 1
            }
        }


class FaultInjectionRequest(BaseModel):
    """Fault injection request"""
    fault: bool
    
    class Config:
        json_schema_extra = {
            "example": {
                "fault": True
            }
        }


# Create FastAPI app and room instance
app = FastAPI(
    title="Room Simulator API",
    description="Simulates a room with two light devices and one brightness sensor",
    version="1.0.0"
)
room = Room()


@app.get("/")
async def root():
    """API root path, returns welcome message"""
    return {
        "message": "Welcome to Room Simulator API",
        "docs": "/docs",
        "status": "running"
    }


@app.get("/status")
async def get_status():
    """Get status of all devices in the room"""
    return room.get_all_status()


@app.get("/lights/{light_name}")
async def get_light(light_name: str):
    """
    Get set state of a specified light (returns plain text: 0 or 1)
    
    Args:
        light_name: Light name (light1 or light2)
    """
    try:
        status = room.get_light_status(light_name)
        return Response(content=str(status["set_state"]), media_type="text/plain")
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@app.put("/lights/{light_name}")
async def set_light(light_name: str, request: Request):
    """
    Set state of a specified light (accepts plain text: 0 or 1)
    
    Args:
        light_name: Light name (light1 or light2)
        request: Request object containing plain text body (0 or 1)
    """
    try:
        state = await request.body()
        state_str = state.decode('utf-8').strip()
        state_int = int(state_str)
        if state_int not in [0, 1]:
            raise HTTPException(
                status_code=400,
                detail="State must be 0 (off) or 1 (on)"
            )
        room.set_light(light_name, state_int)
        return Response(content="ok", media_type="text/plain")
    except ValueError as e:
        if "Unknown light device" in str(e):
            raise HTTPException(status_code=404, detail=str(e))
        raise HTTPException(status_code=400, detail="State must be 0 or 1")


@app.post("/lights/{light_name}/fault")
async def inject_fault(light_name: str, request: Request):
    """
    Inject or remove fault from a specified light

    Args:
        light_name: Light name (light1 or light2)
        request: Request body containing plain text input ('0' or '1')
    """
    try:
        body = await request.body()
        fault_str = body.decode('utf-8').strip()
        if fault_str not in ['0', '1']:
            raise HTTPException(
                status_code=400,
                detail="Fault must be '0' (remove) or '1' (inject)"
            )
        fault = fault_str == '1'
        return room.inject_fault(light_name, fault)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@app.get("/brightness")
async def get_brightness():
    """Get brightness sensor value (returns plain text: 0, 1, or 2)"""
    brightness_info = room.get_brightness()
    return Response(content=str(brightness_info["brightness"]), media_type="text/plain")


def main():
    """Start room simulator service"""
    parser = argparse.ArgumentParser(description="Room Simulator with light devices and brightness sensor")
    parser.add_argument(
        "--port",
        type=int,
        default=8000,
        help="Port to run the server on (default: 8000)"
    )
    parser.add_argument(
        "--host",
        type=str,
        default="0.0.0.0",
        help="Host to run the server on (default: 0.0.0.0)"
    )
    
    args = parser.parse_args()
    
    logger.info("=" * 60)
    logger.info("Room Simulator Service Starting")
    logger.info(f"Server will run on {args.host}:{args.port}")
    logger.info("=" * 60)
    
    # Start FastAPI service
    uvicorn.run(
        app,
        host=args.host,
        port=args.port,
        log_level="info"
    )


if __name__ == "__main__":
    main()
