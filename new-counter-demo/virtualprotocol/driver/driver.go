package driver

import (
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
	"github.com/spf13/cast"
	"k8s.io/klog/v2"
)

func NewClient(protocol ProtocolConfig) (*CustomizedClient, error) {
	client := &CustomizedClient{
		ProtocolConfig: protocol,
		deviceMutex:    sync.Mutex{},
		CounterCli: nil,
	}
	return client, nil
}

func (c *CustomizedClient) InitDevice() error {
	c.CounterCli=NewCounterClient()

	c.CounterCli.InitCounterCli()
	klog.Infoln("Counter Device Init Success")
	return nil
}

func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()
	switch visitor.Name {
		case "count":
			return cast.ToString(c.CounterCli.GetCounter()), nil
		case "status":
			return cast.ToString(c.CounterCli.GetStatus()), nil
		default:
			klog.Infoln("Invalid visitor name")
			return nil, nil
	}
}

func (c *CustomizedClient) DeviceDataWrite(visitor *VisitorConfig, deviceMethodName string, propertyName string, data interface{}) error {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()
	switch deviceMethodName {
	case "UpdateStatus":
		value:= cast.ToInt(data)
		if value==ON||value==OFF||value==PAUSED{
			c.CounterCli.SetStatus(value)
		}else {
			klog.Infoln("Invalid status set:",data)
		}
	case "SetCount":
		value:= cast.ToInt(data)
		if value>=0{
			c.CounterCli.SetCurValue(value)
		}else {
			klog.Infoln("Invalid count set:",data)
		}
	default:
		klog.Infoln("Invalid device method name")
	}
	return nil
}

func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	return nil
}

func (c *CustomizedClient) StopDevice() error {
	c.CounterCli.Close()
	return nil
}

func (c *CustomizedClient) GetDeviceStates() (string, error) {
	return common.DeviceStatusOK, nil
}
