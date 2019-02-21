# Light Mapper
 
 
 ## Description
 
Light Mapper contains code to control an LED light connected to a raspberry Pi through gpio.

<img src="images/raspberry-pi.png">
  
The following diagram has been followed to make the connection with the 
LED in this case :-

<img src="images/raspberry-pi-wiring.png">

Here we are using a push button  switch to test the working condition of the LED, through an independent  circuit.

Depending on the expected state of the light, the program controls whether or not to provide power in pin-18 of the gpio.
When power is provided in the pin, the LED glows (ON State) and when no power is provided on it then it does not glow (OFF state).



## Prerequisites 

### Hardware Prerequisites

1. RaspBerry-Pi (RaspBerry-Pi 3 has been used for this demo)
2. GPIO
3. Breadboard along with wires 
4. LED light
5. Push Button switch (to test the working condition of the light, this can be skipped if needed)

### Software Prerequisites
 
1. Golang (Version 1.11.4 has been used for this demo)
2. KubeEdge (commitID: 4deb29406f470e409a73f2fbbc631705f02899b1 has been used for this demo)
3. Vendor Packages: 
 The following packages need to be present in the GOPATH before building this demo :-
 
           1. github.com/eclipse/paho.mqtt.golang v1.1.1
           2. github.com/stianeikeland/go-rpio v4.4.0
         
                       

## Steps to reproduce

1. Connect the LED to the RaspBerry-Pi using the GPIO as shown in the [circuit diagram](images/raspberry-pi-wiring.png) above.   

2. Clone and run KubeEdge. 
    Please click [KubeEdge Usage](https://github.com/kubeedge/kubeedge#usage) for instructions on the usage of KubeEdge.
Note: Please start KueEdge with in internal or double MQTT mode.

3. Create the LED device in the cloud, with a device twin attribute called "Power Status" (This is used to turn the device ON/OFF). 

4.  Clone the led_raspberrypi_demo 
 
    ```shell
                git clone https://github.com/kubeedge/examples.git $GOPATH/src/github.com/kubeedge/examples
                cd $GOPATH/src/github.com/kubeedge/examples/led_raspberrypi_demo/light_mapper
     ```
 5. Cross compile the mapper to run in RaspBerry-Pi.

    ```shell         
                sudo apt-get install gcc-arm-linux-gnueabi
                export GOARCH=arm
                export GOOS="linux"
                export GOARM=6                             #Pls give the appropriate arm version of your device                               
                export CGO_ENABLED=1
                export CC=arm-linux-gnueabi-gcc
                go build light_mapper.go
    ```
 
 6. Run the mapper in the RaspBerry-Pi while providing the deviceID of the device created and the pin number of the 
  GPIO (18 in this case) as a first and second command line parameter respectively  i.e. 
 "./light_mapper \<deviceID> \<pinNumber>"
        
    ```shell
            Ex :- $ ./light_mapper  2a22d02f-0773-4f88-8c0a-9dccbc50e104 18
     ```
 
  7. Change the device Twin attribute (expected value) "Power State" of the device to "ON" to turn on the light, and 
 "OFF" to turn off the light. The mapper will control the LED to match the state mentioned in the cloud and also report back 
 the actual state of the light to the cloud after updating.

 