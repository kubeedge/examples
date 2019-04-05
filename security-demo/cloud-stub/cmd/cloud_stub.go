package main

import (
	"fmt"
	//"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
	"github.com/kubeedge/examples/security-demo/cloud-stub/cmd/config"
	"github.com/kubeedge/examples/security-demo/cloud-stub/pkg/controller/cloudhub"
)

const (
	STUB_CONFIG_FILE = "config/app.conf"
)

func main() {
	config, err := config.ParseConfig(STUB_CONFIG_FILE)
	if err != nil {
		fmt.Println("Error : failed to parse config with error - %s", err.Error())
	}

	stub := cloudhub.NewCloudStub(config)

	//go stub.Start()
	stub.PlacementServer()
}
