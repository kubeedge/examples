package driver

import (
	"strconv"
	"sync"

	"k8s.io/klog/v2"
)

func NewClient(protocol ProtocolConfig) (*CustomizedClient, error) {
	client := &CustomizedClient{
		ProtocolConfig: protocol,
		deviceMutex:    sync.Mutex{},
		// TODO initialize the variables you added
	}
	return client, nil
}

func (c *CustomizedClient) InitDevice() error {
	// TODO: add init operation
	// you can use c.ProtocolConfig
	return nil
}

func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	// TODO: add the code to get device's data
	// you can use c.ProtocolConfig and visitor
	vLED := GetVLEDStatus()
	return vLED, nil
}

func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	// TODO: set device's data
	// you can use c.ProtocolConfig and visitor

	str, ok := data.(string)
	if !ok {
		klog.Infoln("Error: data is not a string")
	}
	// Convert string to int
	number, err := strconv.Atoi(str)
	if err != nil {
		klog.Infoln("Error converting string to int:", err)
	}
	SetVLEDStatus(number)
	return nil
}

func (c *CustomizedClient) StopDevice() error {
	// TODO: stop device
	// you can use c.ProtocolConfig
	return nil
}
