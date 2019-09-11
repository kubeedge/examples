# Bluetooth With KubeEdge Using CC2650


Users can make use of KubeEdge platform to connect and control their bluetooth devices, provided, the user is aware of the of the data sheet information for their device.
Kubernetes Custom Resource Definition (CRD) and KubeEdge bluetooth mapper is being used to support this feature, using which users can control their device from the cloud. Texas Instruments [CC2650 SensorTag device](http://processors.wiki.ti.com/index.php/CC2650_SensorTag_User%27s_Guide) is being shown here as an example.


## Description

KubeEdge support for bluetooth protocol has been demonstrated here by making use of Texas Instruments CC2650 SensorTag device.
This section contains instructions on how to make use of bluetooth mapper of KubeEdge to control CC2650 SensorTag device.
  
  We will only be focusing on the following features of CC2650 :-
  
  ```shell
  1. IR Temperature
  2. IO-Control :
      2.1 Red Light
      2.2 Greem Light
      2.3 Buzzer
      2.4 Red Light with Buzzer
      2.5 Green Light with Buzzer
      2.6 Red Light along with Green Light
      2.7 Red Light, Green Light along with Buzzer 
      
  ```
  
  The bluetooth mapper has the following major components :-
   - Action Manager
   - Scheduler
   - Watcher
   - Controller
   - Data Converter
  
  More details on bluetooth mapper can be found [here](https://github.com/kubeedge/kubeedge/blob/master/device/bluetooth_mapper/README.md).
  
  
## Prerequisites 

### Hardware Prerequisites

1. Texas instruments CC2650 bluetooth device
2. Linux based edge node with bluetooth support  (An Ubuntu 18.04 laptop has been used in this demo) 

### Software Prerequisites
 
1. Golang (Version 1.11.4 has been used for this demo)
2. KubeEdge (Version 0.3 has been used for this demo)

## Steps to reproduce

1. Clone and run KubeEdge. 
    Please click [KubeEdge Usage](https://github.com/kubeedge/kubeedge/blob/master/docs/getting-started/usage.md) for instructions on the usage of KubeEdge.
    Please ensure that the kubeedge setup is up and running before execution of step 4 (mentioned below).

2. Clone the kubeedge/examples repository.

```console
           git clone https://github.com/kubeedge/examples.git $GOPATH/src/github.com/kubeedge/examples
```

3. Create the CC2650 SensorTag device model and device instance.

```console
           cd $GOPATH/src/github.com/kubeedge/examples/bluetooth-CC2650-demo/sample-crds
           kubectl apply -f CC2650-device-model.yaml
           kubectl apply -f CC2650-device-instance.yaml

           Note: You can change the CRDs to match your requirement
```
 
4. Please ensure that bluetooth service of your device is ON

5. Set 'bluetooth=true' label for the node (This label is a prerequisite for the scheduler to schedule bluetooth_mapper pod on the node [which meets the hardware / software prerequisites] )

```console
kubectl label nodes <name-of-node> bluetooth=true
```

6. Copy the configuration file that has been provided, into its correct path. Please note that the configuration file can be altered as to suit your requirement

```console
cp $GOPATH/src/github.com/kubeedge/examples/bluetooth-CC2650-demo/config.yaml  

$GOPATH/src/github.com/kubeedge/kubeedge/device/bluetooth_mapper/configuration/
``` 

7. Build the mapper by following the steps given below.

```console
cd $GOPATH/src/github.com/kubeedge/kubeedge/device/bluetooth_mapper
make bluetooth_mapper_image
docker tag bluetooth_mapper:v1.0 <your_dockerhub_username>/bluetooth_mapper:v1.0
docker push <your_dockerhub_username>/bluetooth_mapper:v1.0

Note: Before trying to push the docker image to the remote repository please ensure that you have signed into docker from your node, if not please type the followig command to sign in
docker login
# Please enter your username and password when prompted
```

8. Deploy the mapper by following the steps given below.

```console
cd $GOPATH/src/github.com/kubeedge/kubeedge/device/bluetooth_mapper

# Please enter the following details in the deployment.yaml :-
#    1. Replace <edge_node_name> with the name of your edge node at spec.template.spec.voluems.configMap.name
#    2. Replace <your_dockerhub_username> with your dockerhub username at spec.template.spec.containers.image

kubectl create -f deployment.yaml
```

9. Turn ON the CC2650 SensorTag device  

10. The bluetooth mapper is now running, You can monitor the logs of the mapper by using docker logs. You can also play around with the device twin state by altering the desired property in the device instance 
and see the result reflect on the SensorTag device. The configurations of the bluetooth mapper can be altered at runtime Please click [Runtime Configurations](https://github.com/kubeedge/kubeedge/blob/master/device/bluetooth_mapper/README.md#runtime-configuration-modifications) 
for more details.  
