# Temperature Sensor WebUI

## Quick Start
1. Build the image
```
docker build -t temperature-webui:v1.0 .
```
2. Load image to the cloud node
```
// If using kind
kind load docker-image temperature-webui:v1.0 --nodes=<your_cloud_node_name>
```
3. Deploy at the cloud
```
kubectl apply -f ./resource/deploy.yaml
```
4. Open the Web Browser and go to the URL
```
http://<your-cloud-node-ip>:8080
```
