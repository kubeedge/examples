package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type ReadDataFn struct {
	ID         int64  `json:"ID,omitempty"`
	Device     string `json:"Device,omitempty"`
	CreateTime time.Time  `json:"CreateTime,omitempty"`
	Label      string `json:"Label,omitempty"`
	Units      string `json:"Units,omitempty"`
	Value      string `json:"Value,omitempty"`
	MachineID  string `json:"MachineID,omitempty"`
}

func main() {
	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("send-client"),
	})
	if err != nil {
		panic(err)
	}
	Data := ReadDataFn{Device: "cc2650",
		CreateTime: time.Now().Local(),
		Label:      "Temperature",
		Units:      "DegreeCelcius",
		Value:      "85",
		MachineID:  "M1010",
	}

	var i int64
	for i = 0; i < 20; i++ {
		time.Sleep(2 * time.Second)
		Data.ID = i
		Data.CreateTime = time.Now().Local()
		if i == 1 {
			Data.Value = "90"
			Data.MachineID = "M2000"
		} else if i > 2 && i < 5 {
			Data.Value = "100"
			Data.MachineID = "M3030"
		} else if i > 5 && i < 8 {
			Data.Value = "95"
			Data.MachineID = "M5050"
		} else if i == 8 {
			Data.Value = "75"
			Data.MachineID = "M6000"
		} else {
			Data.Value = "80"
			Data.MachineID = "M7000"
		}

		// Publish a message.
		bytes, _ := json.Marshal(Data)
		err = cli.Publish(&client.PublishOptions{
			QoS:       mqtt.QoS0,
			TopicName: []byte("test"),
			Message:   bytes,
		})
		if err != nil {
			panic(err)
		}
	}

	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
