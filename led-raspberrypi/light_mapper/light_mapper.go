package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/kubeedge/examples/led-raspberrypi/light_driver"
)

var (
	DeviceETPrefix            = "$hw/events/device/"
	DeviceETStateUpdateSuffix = "/state/update"
	TwinETUpdateSuffix        = "/twin/update"
	TwinETCloudSyncSuffix     = "/twin/cloud_updated"
	TwinETGetResultSuffix     = "/twin/get/result"
	TwinETGetSuffix           = "/twin/get"
)

var Token_client Token
var ClientOpts *MQTT.ClientOptions
var Client MQTT.Client
var wg sync.WaitGroup
var deviceTwinResult DeviceTwinUpdate

//deviceID contains the device ID of the device provided as a command-line parameter
var deviceID = os.Args[1]

//pinNumber contains the pin number of the GPIO is provided as a command-line parameter
var pinNumber, _ = strconv.ParseInt(os.Args[2], 10, 64)

//Token interface to validate the MQTT connection.
type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}

//DeviceStateUpdate is the structure used in updating the device state
type DeviceStateUpdate struct {
	State string `json:"state,omitempty"`
}

//BaseMessage the base struct of event message
type BaseMessage struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

//TwinValue the struct of twin value
type TwinValue struct {
	Value    *string        `json:"value, omitempty"`
	Metadata *ValueMetadata `json:"metadata,omitempty"`
}

//ValueMetadata the meta of value
type ValueMetadata struct {
	Timestamp int64 `json:"timestamp, omitempty"`
}

//TypeMetadata the meta of value type
type TypeMetadata struct {
	Type string `json:"type,omitempty"`
}

//TwinVersion twin version
type TwinVersion struct {
	CloudVersion int64 `json:"cloud"`
	EdgeVersion  int64 `json:"edge"`
}

//MsgTwin the struct of device twin
type MsgTwin struct {
	Expected        *TwinValue    `json:"expected,omitempty"`
	Actual          *TwinValue    `json:"actual,omitempty"`
	Optional        *bool         `json:"optional,omitempty"`
	Metadata        *TypeMetadata `json:"metadata,omitempty"`
	ExpectedVersion *TwinVersion  `json:"expected_version,omitempty"`
	ActualVersion   *TwinVersion  `json:"actual_version,omitempty"`
}

//DeviceTwinUpdate the struct of device twin update
type DeviceTwinUpdate struct {
	BaseMessage
	Twin map[string]*MsgTwin `json:"twin"`
}

// Initiate the MQTT connection
func init() {
	ClientOpts = HubClientInit("tcp://127.0.0.1:1884", "eventbus", "", "")
	Client = MQTT.NewClient(ClientOpts)
	if Token_client = Client.Connect(); Token_client.Wait() && Token_client.Error() != nil {
		fmt.Println("client.Connect() Error is ", Token_client.Error())
	}
}

// HubclientInit create mqtt client config
func HubClientInit(server, clientID, username, password string) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientID).SetCleanSession(true)
	if username != "" {
		opts.SetUsername(username)
		if password != "" {
			opts.SetPassword(password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)
	return opts
}

//changeDeviceState function is used to change the state of the device
func changeDeviceState(state string) {
	fmt.Println("Changing the state of the device to online")
	var deviceStateUpdateMessage DeviceStateUpdate
	deviceStateUpdateMessage.State = state
	stateUpdateBody, err := json.Marshal(deviceStateUpdateMessage)
	if err != nil {
		fmt.Println("Error:   ", err)
	}
	deviceStatusUpdate := DeviceETPrefix + deviceID + DeviceETStateUpdateSuffix
	Token_client = Client.Publish(deviceStatusUpdate, 0, false, stateUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		fmt.Println("client.publish() Error in device state update  is ", Token_client.Error())
	}
}

//changeTwinValue sends the updated twin value to the edge through the MQTT broker
func changeTwinValue(updateMessage DeviceTwinUpdate) {
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		fmt.Println("Error:   ", err)
	}
	deviceTwinUpdate := DeviceETPrefix + deviceID + TwinETUpdateSuffix
	Token_client = Client.Publish(deviceTwinUpdate, 0, false, twinUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		fmt.Println("client.publish() Error in device twin update is ", Token_client.Error())
	}
}

//syncToCloud function syncs the updated device twin information to the cloud
func syncToCloud(updateMessage DeviceTwinUpdate) {
	deviceTwinResultUpdate := DeviceETPrefix + deviceID + TwinETCloudSyncSuffix
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		fmt.Println("Error:   ", err)
	}
	Token_client = Client.Publish(deviceTwinResultUpdate, 0, false, twinUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		fmt.Println("client.publish() Error in device twin update is ", Token_client.Error())
	}
}

// OnSubMessageReceived callback function which is called when message is received
func OnSubMessageReceived(client MQTT.Client, message MQTT.Message) {
	err := json.Unmarshal(message.Payload(), &deviceTwinResult)
	if err != nil {
		fmt.Println("Error in unmarshalling:  ", err)
	}
}

//createActualUpdateMessage function is used to create the device twin update message
func createActualUpdateMessage(actualValue string) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	actualMap := map[string]*MsgTwin{"Power_Status": {Actual: &TwinValue{Value: &actualValue}, Metadata: &TypeMetadata{Type: "Updated"}}}
	deviceTwinUpdateMessage.Twin = actualMap
	return deviceTwinUpdateMessage
}

//getTwin function is used to get the device twin details from the edge
func getTwin(updateMessage DeviceTwinUpdate) {
	getTwin := DeviceETPrefix + deviceID + TwinETGetSuffix
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		fmt.Println("Error:   ", err)
	}
	Token_client = Client.Publish(getTwin, 0, false, twinUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		fmt.Println("client.publish() Error in device twin get  is ", Token_client.Error())
	}
}

//subscribe function subscribes  the device twin information through the MQTT broker
func subscribe() {
	for {
		getTwinResult := DeviceETPrefix + deviceID + TwinETGetResultSuffix
		Token_client = Client.Subscribe(getTwinResult, 0, OnSubMessageReceived)
		if Token_client.Wait() && Token_client.Error() != nil {
			fmt.Println("subscribe() Error in device twin result get  is ", Token_client.Error())
		}
		time.Sleep(1 * time.Second)
		if deviceTwinResult.Twin != nil {
			wg.Done()
			break
		}
	}
}

//equateTwinValue is responsible for equating the actual state of the device to the expected state that has been set
func equateTwinValue(updateMessage DeviceTwinUpdate) {
	fmt.Println("Watching on the device twin values for device with deviceID: ", os.Args[1])
	wg.Add(1)
	go subscribe()
	getTwin(updateMessage)
	wg.Wait()
	if deviceTwinResult.Twin["Power_Status"].Actual == nil || *deviceTwinResult.Twin["Power_Status"].Expected.Value != *deviceTwinResult.Twin["Power_Status"].Actual.Value {
		fmt.Println("Expected Value : ", *deviceTwinResult.Twin["Power_Status"].Expected.Value)
		if deviceTwinResult.Twin["Power_Status"].Actual == nil {
			fmt.Println("Actual Value: ", deviceTwinResult.Twin["Power_Status"].Actual)
		} else {
			fmt.Println("Actual Value: ", *deviceTwinResult.Twin["Power_Status"].Actual.Value)
		}
		fmt.Println("Equating the actual  value to expected value")
		switch strings.ToUpper(*deviceTwinResult.Twin["Power_Status"].Expected.Value) {
		case "ON":
			fmt.Println("Turning ON the light")
			//Turn On the light by supplying power on the pin specified
			lightdriver.TurnON(pinNumber)

		case "OFF":
			fmt.Println("Turning OFF the light")
			//Turn Off the light by cutting off power on the pin specified
			lightdriver.TurnOff(pinNumber)

		default:
			panic("OOPS!!!!! Attempt to perform invalid operation " + *deviceTwinResult.Twin["Power_Status"].Expected.Value + " on LED light")
		}
		updateMessage = createActualUpdateMessage(*deviceTwinResult.Twin["Power_Status"].Expected.Value)
		changeTwinValue(updateMessage)
		time.Sleep(2 * time.Second)
		fmt.Println("Syncing to cloud.....")
		syncToCloud(updateMessage)
	} else {
		fmt.Println("Actual values are in sync with Expected value")
	}
}

func main() {
	changeDeviceState("online")
	updateMessage := createActualUpdateMessage("unknown")
	for {
		equateTwinValue(updateMessage)
	}
}
