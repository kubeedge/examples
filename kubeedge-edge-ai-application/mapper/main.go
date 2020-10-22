package main

import (
	"encoding/json"
	"fmt"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"strings"
	"sync"
)

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
type deviceState struct {
	Consumption int    `json:"consumption"`
	Linkquality int    `json:"linkquality"`
	Power       int    `json:"power"`
	State       string `json:"state"`
	Temperature int    `json:"temperature"`
}

var wg sync.WaitGroup

func main() {
	// connect to mqtt broker
	cli := connectToMqtt()
	// make goroutine a receive face recong result
	data := make(chan string)
	// get command from face recon
	go subRes(cli, "faceRecon", data)
	wg.Add(1)
	// send command to device
	go publicTodevice(cli, data)
	wg.Add(1)
	sta := make(chan string)
	//get state from device
	go subRes(cli, "zigbee2mqtt/0x00158d0003933562", sta)
	wg.Add(1)
	//send state to cloud
	go publishToMqtt(cli, "$hw/events/device/switch/twin/update", sta)
	wg.Add(1)
	wg.Wait()
}

//cli sub result topic 接收人脸识别数据 当为正常的时候打开开关 向zigbee2mqtt发送数据 协程a
//设定一个全局变量，表征是否要打开开关 全局变量
//cli sub zigbee2mqtt topic 当有数据时 同步至云端 协程b
func subRes(cli *client.Client, topic string, data chan string) {
	err := cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(topic),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {
					fmt.Println(string(topicName), string(message))
					data <- string(message)
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	wg.Done()
}
func publicTodevice(cli *client.Client, data <-chan string) {
	topic := "zigbee2mqtt/0x00158d0003933562/set"
	for {
		select {
		case state := <-data:
			state = strings.Replace(state, "\n", "", -1)
			actualMap := map[string]string{"state": state}
			body, _ := json.Marshal(actualMap)
			cli.Publish(&client.PublishOptions{
				TopicName: []byte(topic),
				QoS:       mqtt.QoS0,
				Message:   body,
			})
		default:
			continue
		}
	}

	wg.Done()
}
func connectToMqtt() *client.Client {
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})
	defer cli.Terminate()
	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "10.28.160.123:1883",
		ClientID: []byte("receive-client"),
	})
	if err != nil {
		panic(err)
	}
	return cli
}
func publishToMqtt(cli *client.Client, topic string, sta chan string) {
	s := new(deviceState)
	for {
		select {
		case jso := <-sta:
			err := json.Unmarshal([]byte(jso), &s)
			if err != nil {
				fmt.Println("error")
				panic(err)
			}
			fmt.Println("state:" + s.State)
			updateMessage := createActualUpdateMessage(s.State)
			twinUpdateBody, _ := json.Marshal(updateMessage)

			cli.Publish(&client.PublishOptions{
				TopicName: []byte(topic),
				QoS:       mqtt.QoS0,
				Message:   twinUpdateBody,
			})
		default:
			continue
		}
	}

}

//createActualUpdateMessage function is used to create the device twin update message
func createActualUpdateMessage(actualValue string) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	actualMap := map[string]*MsgTwin{"state": {Actual: &TwinValue{Value: &actualValue}, Metadata: &TypeMetadata{Type: "Updated"}}}
	deviceTwinUpdateMessage.Twin = actualMap
	return deviceTwinUpdateMessage
}
