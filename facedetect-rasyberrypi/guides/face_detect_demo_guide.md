# Guide for use Kubeedge to run face detect Demo for edge at a raspberryPi

## 1. Prepare Kubeedge environment 
install guide please refers to the kubeedge install readme
[Kubeedge installation reference](https://github.com/kubeedge/kubeedge/blob/master/README.md)

## 2. Prepare the raspberryPi environment
> Demo supports raspberrypi3 and higher
```shell script
run rasp-config to enable the raspberryPi camera. and reboot.
```
## 3 . Run the face detect demo 
> Provide two ways to run this face recognition demo on RaspberryPi

### 3.1. By using rtsp stream provided by the raspberryPi built-in camera
> note: This method supports the camera which support RTSP video stream. Not just the camera that comes with raspberryPi \
> The following steps use Raspberry Pi as an example

+ Step1. Generate RTSP vide stream with build-in camera

> install the dependent packages
```shell script
apt-get install cmake liblog4cpp5-dev libv4l-dev
```
> [download the rtsp tool](https://github.com/mpromonet/v4l2rtspserver/releases/tag/v0.1.9)
```shell script
# install the rtsp tool
dpkg -i v4l2rtspserver-0.1.9-Linux-armv7l.deb
```
> start the RTSP server
```shell script
v4l2rtspserver -H 640 -W 480 -F 15 -P 8555 /dev/video0 > rtsp.log 2>&1&
```
> -F : the frame rate \
> -H : the frame height \
> -W : the frame width \
> -P : the rtsp server port
> after finished the above steps, you will get a rtsp server url
```shell script
rtsp://hostip:8555/unicast
# the hostip depends on the configuration on your machine
```

+ Step2. Pull the Demo images
```shell script
docker pull kubeedge/face-detect-demo:v1
```
+ Step3. Run with docker
```shell script
docker run -d -p 8099:8099 -e VIDEO_URL=rtsp://${hostip}:8555/unicast kubeedge/face-detect-demo:v1 python3 /opt/src/app_rtsp.py
```
+ Or run with kubeedge
> The deployment files are as follows \
> Please modify the ${hostip} as your machine's \
> Please modify the ${node_id} maybe you need replace the nodeSelector \
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: face-detect-demo
  labels:
    app: demo
spec:
  containers:
  - name: demo
    image: kubeedge/face-detect-demo:v1
    imagePullPolicy: IfNotPresent
    ports:
    - containerPort: 8099
      hostPort: 8099
    args:
    - /opt/src/app_rtsp.py
    command:
    - python3
    env:
    - name: VIDEO_URL
      value: rtsp://${hostip}:8555/unicast
  nodeSelector:
    name: ${node_id}
```


### 3.2. By using the picamera module in python (Recommended,Performance will be better)
> this method just support the raspberryPi, does not support other types of machine
+ Step1. Pull the Demo images
```shell script
docker pull docker pull kubeedge/face-detect-demo:v1
```

+ Step2. Run with docker
> The deployment files are as follows
```shell script
docker run -d --privileged=true -v /dev/:/dev/ -p 8099:8099 kubeedge/face-detect-demo:v1 python3 /opt/src/app.py
```

+ Or run with kubeedge
> The deployment files are as follows \
> Please modify the ${node_id} maybe you need replace the nodeSelector
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: face-detect-demo
  labels:
    app: demo
spec:
  containers:
  - name: demo
    image: kubeedge/face-detect-demo:v1
    imagePullPolicy: IfNotPresent
    ports:
    - containerPort: 8099
      hostPort: 8099
    args:
    - /opt/src/app.py
    command:
    - python3
    volumeMounts:
    - name: dev
      mountPath: /dev/
    securityContext:
      privileged: true
  volumes:
  - name: dev
    hostPath:
      path: /dev/
  nodeSelector:
    name: ${node_id}
```
