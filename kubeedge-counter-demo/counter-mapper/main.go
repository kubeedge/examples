package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	counter "github.com/kubeedge/examples/kubeedge-counter-demo/counter-mapper/device"
	"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/types"
)

var cli mqtt.Client
var counterList map[string]*counter.Counter

const (
	mqttUrl  = "tcp://127.0.0.1:1883"
	subTopic = "$hw/events/device/+/twin/update"
	pubTopic = "$hw/events/device/%s/twin/update"
)

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

//createActualUpdateMessage function is used to create the device twin update message
func createActualUpdateMessage(actualValue string) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	actualMap := map[string]*MsgTwin{"status": {Actual: &TwinValue{Value: &actualValue}, Metadata: &TypeMetadata{Type: "Updated"}}}
	deviceTwinUpdateMessage.Twin = actualMap
	return deviceTwinUpdateMessage
}

func publishToMqtt(name string, data int) {
	updateMessage := createActualUpdateMessage(strconv.Itoa(data))
	twinUpdateBody, _ := json.Marshal(updateMessage)

	topic := fmt.Sprintf(pubTopic, name)
	token := cli.Publish(topic, 0, false, twinUpdateBody)

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func connectToMqtt() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttUrl)

	cli = mqtt.NewClient(opts)

	token := cli.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	return cli
}

func main() {
	stopchan := make(chan os.Signal)
	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGKILL)
	defer close(stopchan)

	cli = connectToMqtt()
	counterList = make(map[string]*counter.Counter)

	token := cli.Subscribe(subTopic+"/document", 0, func(client mqtt.Client, msg mqtt.Message) {
		Update := &types.DeviceTwinDocument{}
		err := json.Unmarshal(msg.Payload(), Update)
		if err != nil {
			fmt.Printf("Unmarshal error: %v\n", err)
		}

		ctr := GetCounterByName(GetNameFromTopic(msg.Topic()))

		cmd := *Update.Twin["status"].CurrentState.Expected.Value

		if cmd == "ON" && cmd != ctr.CurrentStatus {
			ctr.TurnOn()
			fmt.Printf("turn on counter %s.\n", ctr.Name)
		}

		if cmd == "OFF" && cmd != ctr.CurrentStatus {
			ctr.TurnOff()
			fmt.Printf("turn off counter %s.\n", ctr.Name)
		}

		ctr.CurrentStatus = cmd
	})

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	select {
	case <-stopchan:
		fmt.Printf("Interrupt, exit.\n")
		break
	}
}

func GetCounterByName(name string) *counter.Counter {
	c, ok := counterList[name]
	if ok {
		return c
	}

	ctr := counter.NewCounter(name, publishToMqtt)
	counterList[name] = ctr
	return ctr
}

func GetNameFromTopic(topic string) string {
	ts := strings.Split(topic, "/")
	return ts[3]
}
