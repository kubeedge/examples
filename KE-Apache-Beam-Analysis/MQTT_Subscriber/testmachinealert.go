package main

import (
	"fmt"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func main() {
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})
	stopchan := make(chan int)

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "127.0.0.1:1883",
		ClientID: []byte("receive-client"),
	})
	if err != nil {
		panic(err)
	}
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("alert"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {

					str := string(message)
					fmt.Println("", str)
				},
			},
		},
	})
	<-stopchan
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection Successful:")
	}
}
