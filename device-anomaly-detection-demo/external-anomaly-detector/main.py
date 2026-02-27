import grpc
from concurrent import futures
import time
import argparse
import threading
import logging
import os
import sys
from kubernetes import client, config

# Add externalanomalydetectorpb to sys.path to fix generated code imports
sys.path.append(os.path.join(os.path.dirname(__file__), 'externalanomalydetectorpb'))

from externalanomalydetectorpb import external_anomaly_detector_pb2
from externalanomalydetectorpb import external_anomaly_detector_pb2_grpc
from model_training.model_training import predict_brightness

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Data store
data_store = {
    "light-1": None,
    "light-2": None,
    "brightness-sensor-1": None
}
data_lock = threading.Lock()

class ExternalAnomalyDetector(external_anomaly_detector_pb2_grpc.ExternalAnomalyDetectorServicer):
    def DetectAnomaly(self, request, context):
        logger.info(f"Received request for namespace={request.namespace}, device={request.deviceName}, property={request.propertyName}")
        
        try:
            # Filter and store interested devices data in memory
            if request.namespace == 'default':
                if request.deviceName in ['light-1', 'light-2'] and request.propertyName == 'working_state':
                    # Store received device data in memory
                    self._store_device_data(request.deviceName, request.propertyName, request.propertyValue)
                elif request.deviceName == 'brightness-sensor-1' and request.propertyName == 'brightness':
                    # Store received device data in memory
                    self._store_device_data(request.deviceName, request.propertyName, request.propertyValue)
            
            return external_anomaly_detector_pb2.AnomalyResponse(
                success=True,
                errorCode=0,
                message="Processed"
            )
        except Exception as e:
            logger.error(f"Error in DetectAnomaly: {e}")
            return external_anomaly_detector_pb2.AnomalyResponse(
                success=False,
                errorCode=1,
                message=str(e)
            )

    def _store_device_data(self, device_name, property_name, property_value):
        """store device data in memory"""
        with data_lock:
            data_store[device_name] = int(property_value)
        logger.info(f"Stored {device_name}.{property_name} = {property_value} in memory")
        
        
def locate_fault(l1, l2, actual_brightness, predicted):
    """Attempt to localize faulty device(s) by flipping light states and re-running Bayesian inference.
    
    Returns a dict mapping each device name to its anomaly message.
    """
    flipped_l1 = 1 - l1
    flipped_l2 = 1 - l2

    # Try flipping light-1 only
    pred = predict_brightness(flipped_l1, l2)
    if pred is not None and int(pred) == actual_brightness:
        logger.info(f"Fault localized: light-1 is faulty (state reported={l1}, should be={flipped_l1})")
        return {
            "light-1": f"Device fault: reported state={l1} but should be={flipped_l1}",
            "light-2": "no anomaly",
            "brightness-sensor-1": "no anomaly",
        }

    # Try flipping light-2 only
    pred = predict_brightness(l1, flipped_l2)
    if pred is not None and int(pred) == actual_brightness:
        logger.info(f"Fault localized: light-2 is faulty (state reported={l2}, should be={flipped_l2})")
        return {
            "light-1": "no anomaly",
            "light-2": f"Device fault: reported state={l2} but should be={flipped_l2}",
            "brightness-sensor-1": "no anomaly",
        }

    # Try flipping both light-1 and light-2
    pred = predict_brightness(flipped_l1, flipped_l2)
    if pred is not None and int(pred) == actual_brightness:
        logger.info(f"Fault localized: both light-1 and light-2 are faulty")
        return {
            "light-1": f"Device fault: reported state={l1} but should be={flipped_l1}",
            "light-2": f"Device fault: reported state={l2} but should be={flipped_l2}",
            "brightness-sensor-1": "no anomaly",
        }

    # Cannot localize fault
    logger.warning("Could not localize fault after trying all flip combinations")
    return {
        "light-1": "Device fault: unknown cause",
        "light-2": "Device fault: unknown cause",
        "brightness-sensor-1": f"Anomaly: predicted brightness={predicted} does not match actual={actual_brightness}",
    }


def anomaly_detection_loop():
    """Background thread that performs anomaly detection every 10 seconds"""
    while True:
        time.sleep(10)
        try:
            # Read data from memory
            with data_lock:
                l1 = data_store.get("light-1")
                l2 = data_store.get("light-2")
                actual_brightness = data_store.get("brightness-sensor-1")
            
            if l1 is not None and l2 is not None and actual_brightness is not None:
                # Predict expected brightness using Bayesian Network model
                predicted = predict_brightness(l1, l2)
                
                if predicted is not None:
                    predicted = int(predicted)
                    
                    logger.info(f"Anomaly detection: light-1={l1}, light-2={l2} -> predicted={predicted}, actual={actual_brightness}")
                    
                    if predicted != actual_brightness:
                        # Anomaly detected — try to localize the faulty device(s)
                        messages = locate_fault(l1, l2, actual_brightness, predicted)
                    else:
                        messages = {
                            "light-1": "no anomaly",
                            "light-2": "no anomaly",
                            "brightness-sensor-1": "no anomaly",
                        }
                    
                    # Update anomaly status to apiserver
                    update_anomaly_status(messages)
                else:
                    logger.warning("Prediction returned None")
            else:
                logger.debug(f"Waiting for data. Current store: {data_store}")
                
        except Exception as e:
            logger.error(f"Error in anomaly detection loop: {e}")

def update_anomaly_status(messages):
    """Update anomaly status for each device to k8s apiserver.
    
    Args:
        messages: dict mapping device name to its anomaly message string.
    """
    try:
        api = client.CustomObjectsApi()
        
        for d_name, message in messages.items():
            patch = {
                "status": {
                    "extensions": {
                        "anomaly": message
                    }
                }
            }
            api.patch_namespaced_custom_object(
                group="devices.kubeedge.io",
                version="v1beta1",
                namespace="default",
                plural="devicestatuses",
                name=d_name,
                body=patch
            )
            logger.info(f"Updated {d_name} anomaly status: {message}")
        
    except Exception as e:
        logger.error(f"Error updating anomaly status to apiserver: {e}")

def serve(host, port):
    """Start gRPC server and anomaly detection loop"""
    # Initialize kubernetes client once
    try:
        config.load_kube_config()
        logger.info("Kubernetes config loaded successfully")
    except Exception as e:
        logger.error(f"Failed to load kubernetes config: {e}")
        sys.exit(1)
    
    # Start gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    external_anomaly_detector_pb2_grpc.add_ExternalAnomalyDetectorServicer_to_server(
        ExternalAnomalyDetector(), server
    )
    server.add_insecure_port(f'{host}:{port}')
    server.start()
    logger.info(f"gRPC server started on {host}:{port}")
    
    # Start background anomaly detection thread
    anomaly_thread = threading.Thread(target=anomaly_detection_loop, daemon=True)
    anomaly_thread.start()
    logger.info("Anomaly detection loop started (runs every 10 seconds)")
    
    server.wait_for_termination()

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='External Anomaly Detector gRPC Server')
    parser.add_argument('--host', default='0.0.0.0', help='Host to bind (default: 0.0.0.0)')
    parser.add_argument('--port', default=6804, type=int, help='Port to bind (default: 6804)')
    args = parser.parse_args()
    
    serve(args.host, args.port)
