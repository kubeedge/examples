package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	// UNCHANGED: this is the counter device logic, no kubeedge dependency
	"github.com/kubeedge/examples/kubeedge-counter-demo/counter-mapper/device"

	// REMOVED: "github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/types"
	// Replaced with local struct definitions below
)

var cli mqtt.Client

const (
	mqttUrl = "tcp://127.0.0.1:1883"
	// topic to PUBLISH actual/reported value to EdgeCore
	topic = "$hw/events/device/counter/twin/update"
)

type BaseMessage struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

type TwinValue struct {
	Value    *string        `json:"value,omitempty"`
	Metadata *ValueMetadata `json:"metadata,omitempty"`
}

type ValueMetadata struct {
	Timestamp int64 `json:"timestamp,omitempty"`
}

type TypeMetadata struct {
	Type string `json:"type,omitempty"`
}

type TwinVersion struct {
	CloudVersion int64 `json:"cloud"`
	EdgeVersion  int64 `json:"edge"`
}

type MsgTwin struct {
	Expected        *TwinValue    `json:"expected,omitempty"`
	Actual          *TwinValue    `json:"actual,omitempty"`
	Optional        *bool         `json:"optional,omitempty"`
	Metadata        *TypeMetadata `json:"metadata,omitempty"`
	ExpectedVersion *TwinVersion  `json:"expected_version,omitempty"`
	ActualVersion   *TwinVersion  `json:"actual_version,omitempty"`
}

type DeviceTwinUpdate struct {
	BaseMessage
	Twin map[string]*MsgTwin `json:"twin"`
}

// --- ADDED: replaces types.DeviceTwinDocument from the removed package ---
// These mirror the document message KubeEdge EventBus publishes on
// $hw/events/device/{id}/twin/update/document

type TwinPropertyValue struct {
	Value    *string        `json:"value,omitempty"`
	Metadata *ValueMetadata `json:"metadata,omitempty"`
}

type TwinState struct {
	Expected *TwinPropertyValue `json:"expected,omitempty"`
	Actual   *TwinPropertyValue `json:"actual,omitempty"`
}

type MsgTwinDoc struct {
	LastState    *TwinState `json:"last,omitempty"`
	CurrentState *TwinState `json:"current,omitempty"`
}

// DeviceTwinDocument is the message published by EdgeCore on the /document topic
// when a twin desired state changes (triggered by cloud patching the Device CR).
type DeviceTwinDocument struct {
	BaseMessage
	Twin map[string]*MsgTwinDoc `json:"twin"`
}

// createActualUpdateMessage builds the MQTT message to report the actual
// counter value back to EdgeCore. UNCHANGED in logic and format.
func createActualUpdateMessage(actualValue string) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	actualMap := map[string]*MsgTwin{
		"status": {
			Actual:   &TwinValue{Value: &actualValue},
			Metadata: &TypeMetadata{Type: "Updated"},
		},
	}
	deviceTwinUpdateMessage.Twin = actualMap
	log.Println("Create actually device twin update...")
	return deviceTwinUpdateMessage
}

// publishToMqtt is called by the counter every second with the current count.
// It publishes the count as the "actual" reported value of the status twin.
// UNCHANGED in logic — EdgeCore picks this up and syncs to cloud CR status.
func publishToMqtt(data int) {
	updateMessage := createActualUpdateMessage(strconv.Itoa(data))
	twinUpdateBody, _ := json.Marshal(updateMessage)
	token := cli.Publish(topic, 0, false, twinUpdateBody)
	log.Println("Publish to mqtt...")
	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
}

// connectToMqtt is UNCHANGED
func connectToMqtt() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttUrl)
	cli = mqtt.NewClient(opts)
	token := cli.Connect()
	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
	return cli
}

func main() {
	stopchan := make(chan os.Signal, 1)
	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGKILL)
	defer close(stopchan)

	cli = connectToMqtt()

	// Wire counter: every tick publishes the count as reported value
	ctr := counter.NewCounter(publishToMqtt)
	current_status := "OFF"
        ctr.TurnOn() 
	log.Println("Start to subscribe the topic document....")

	// Subscribe to the document topic — EdgeCore publishes here whenever
	// the cloud patches Device CR desired state.
	// CHANGED: uses local DeviceTwinDocument instead of types.DeviceTwinDocument
	token := cli.Subscribe(topic+"/document", 0, func(client mqtt.Client, msg mqtt.Message) {
		Update := &DeviceTwinDocument{}
		err := json.Unmarshal(msg.Payload(), Update)
		if err != nil {
			log.Printf("Unmarshal error: %v\n", err)
			return
		}

		// CHANGED: added nil guards — CurrentState can be absent on first message
		twinEntry, ok := Update.Twin["status"]
		if !ok || twinEntry == nil {
			log.Println("No status twin in document, skipping")
			return
		}
		if twinEntry.CurrentState == nil ||
			twinEntry.CurrentState.Expected == nil ||
			twinEntry.CurrentState.Expected.Value == nil {
			log.Println("Incomplete twin state in document, skipping")
			return
		}

		cmd := *twinEntry.CurrentState.Expected.Value
		if cmd == "ON" && cmd != current_status {
			ctr.TurnOn()
			log.Printf("turn on counter.\n")
		}
		if cmd == "OFF" && cmd != current_status {
			ctr.TurnOff()
			log.Printf("turn off counter.\n")
		}
		current_status = cmd
	})

	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	select {
	case <-stopchan:
		log.Printf("Interrupt, exit.\n")
	}
}
