package main

import (
	"flag"
	"fmt"
	"time"
	"encoding/json"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"

	"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/types"
)

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
	optional := false
	songname := flag.String("song", "song1", "Song to be Played")
	flag.Parse()

	data := types.DeviceTwinUpdate{}
	data.EventID = "12345"
	data.Timestamp = time.Now().Unix()/1e6
	data.Twin = map[string]*types.MsgTwin{"track": &types.MsgTwin{Expected:&types.TwinValue{Value:songname, Metadata:&types.ValueMetadata{Timestamp: time.Now().Unix()/1e6}}, Optional: &optional, Metadata:&types.TypeMetadata{Type:"string"}, ExpectedVersion: &types.TwinVersion{CloudVersion: 22022, EdgeVersion:0}}}

	// Publish a message.
	bytes,_:=json.Marshal(data)
	fmt.Println(string(bytes))
	err = cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte("$hw/events/device/speaker-01/twin/cloud_updated"),
		Message:   bytes,
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(2*time.Second)
	// Disconnect the Network Connection.
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
