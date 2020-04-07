package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	parseFlag()

	loadConfigMap()

	InitCLient()

	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("connect info: ", token.Error())
		os.Exit(1)
	}

	//ChangeDeviceState("online")

	//GetTwin(CreateActualDeviceStatus(UNKNOW, UNKNOW, UNKNOW))

	go UpdateActualDeviceStatus()

	for {
		time.Sleep(time.Second * 2)
	}

	Client.Disconnect(250)
}
