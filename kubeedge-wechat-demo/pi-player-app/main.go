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

var musicDir = "/home/pi/music/"

func main() {
	// Get files
	files, err := os.Open(musicDir)
	if err != nil {
		panic(err)
	}
	var m = make(map[string]string)
	list, _ := files.Readdirnames(0)
	fmt.Printf("Files read dirnames result: %v\n", list)
	for _, track := range list {
		trackWithoutSuffix := strings.TrimSuffix(track, ".mp3")
		fmt.Printf("Loading track key: %s value: %s\n", trackWithoutSuffix, track)
		m[trackWithoutSuffix] = track
	}
	fmt.Println("Get music list successfully")

	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	fmt.Println("Create mqtt client successfully")

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
	fmt.Println("Connect mqtt client successfully")

	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			{
				TopicFilter: []byte(`$hw/events/device/speaker-02/twin/update/document`),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					Update := &types.DeviceTwinDocument{}
					err := json.Unmarshal(message, Update)
					if err != nil {
						fmt.Println("Unmarshal error", err)
						fmt.Printf("Unmarshal error: %v\n", err)
					}
					cmd := exec.Command("pkill", "-9", "omxplayer")
					cmd.Run()

					trackToPlay := *Update.Twin["track"].CurrentState.Expected.Value
					fmt.Printf("Receive expected track: %s\n", trackToPlay)

					// Stop music
					if trackToPlay == "stop" {
						return
					}

					_, ok := m[trackToPlay]
					if !ok {
						fmt.Printf("Could not find song %s in playlist\n", trackToPlay)
						trackToPlay = MapRandomKeyGet(m).(string)
						fmt.Printf("Selected random track %s to play\n", trackToPlay)
					}
					fmt.Printf("Playing track: %s\n", musicDir+trackToPlay+".mp3")
					cmd = exec.Command("omxplayer", "-o", "local", musicDir+trackToPlay+".mp3")
					err = cmd.Run()
					if err != nil {
						fmt.Printf("Error while playing track: %v\n", err)
					}
				},
			},
		},
	})
	fmt.Println("Subscribe mqtt topic successfully")

	<-stopchan
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection successfully")
	}
}

func MapRandomKeyGet(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()

	return keys[rand.Intn(len(keys))].Interface()
}
