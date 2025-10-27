package driver

import (
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
)

// CustomizedDev is the customized device configuration and client information.
type CustomizedDev struct {
	Instance         common.DeviceInstance
	CustomizedClient *CustomizedClient

}

type CustomizedClient struct {
	// TODO add some variables to help you better implement device drivers
	deviceMutex sync.Mutex
	ProtocolConfig
	CounterCli CounterClient
}

type ProtocolConfig struct {
	ProtocolName string `json:"protocolName"`
	ConfigData   `json:"configData"`
}

type ConfigData struct {
	// TODO: add your protocol config data
}

type VisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	DataType string `json:"dataType"`
	Name 	string `json:"name"`
}
