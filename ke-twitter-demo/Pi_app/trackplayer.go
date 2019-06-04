package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"

	"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/types"
)

func main() {

	files, err := os.Open("/home/pi/Music/")
	if err != nil {
		panic(err)
	}
	var m = make(map[string]int)
	list, _ := files.Readdirnames(0)
	for _, track := range list {
		trackWithoutSuffix := strings.TrimSuffix(track, ".mp3")
		fmt.Println(trackWithoutSuffix)
		m[trackWithoutSuffix] = 1
	}

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
	err = cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("receive-client"),
	})
	if err != nil {
		panic(err)
	}

	err = cli.Subscribe(&client.SubscribeOptions{

		SubReqs: []*client.SubReq{
			{
				TopicFilter: []byte(`$hw/events/device/speaker-01/twin/update/document`),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					Update := &types.DeviceTwinDocument{}
					err := json.Unmarshal(message, Update)
					if err != nil {
						fmt.Println("error =", err)
					}
					cmd := exec.Command("pkill", "-9", "omxplayer")
					cmd.Run()

					trackToPlay := *Update.Twin["track"].CurrentState.Expected.Value
					_, ok := m[trackToPlay]
					if !ok {
						fmt.Printf("Could not find song %s in playlist\n", trackToPlay)
						trackToPlay = MapRandomKeyGet(m).(string)
						fmt.Printf("Selected random track %s to play\n", trackToPlay)
					}
					fmt.Printf("Playing track : %s\n", "/home/pi/Music/"+trackToPlay+".mp3")
					cmd = exec.Command("omxplayer", "-o", "local", "/home/pi/Music/"+trackToPlay+".mp3")
					err = cmd.Run()
					if err != nil {
						fmt.Printf("error while playing track = %v\n", err)
					}
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

func MapRandomKeyGet(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()

	return keys[rand.Intn(len(keys))].Interface()
}
