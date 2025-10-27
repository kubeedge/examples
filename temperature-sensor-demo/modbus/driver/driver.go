package driver

import (
	"encoding/binary"
	"fmt"
	"github.com/spf13/cast"
	"k8s.io/klog/v2"
	"math"
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
)

func NewClient(protocol ProtocolConfig) (*CustomizedClient, error) {
	client := &CustomizedClient{
		ProtocolConfig: protocol,
		deviceMutex:    sync.Mutex{},
		ModbusClient:   nil,
	}
	return client, nil
}

// InitDevice:init modbus client
func (c *CustomizedClient) InitDevice() error {
	var config interface{}
	klog.Infoln("Modbus CommunicateMode:", c.ProtocolConfig)
	switch c.ProtocolConfig.ConfigData.CommunicateMode {
	case "TCP":
		config = ModbusTCP{
			SlaveID:  byte(c.ConfigData.SlaveID),
			DeviceIP: c.ConfigData.IP,
			TCPPort:  c.ConfigData.Port,
		}
	case "RTU":
		config = ModbusRTU{
			SlaveID:      byte(c.ConfigData.SlaveID),
			SerialName:   c.ConfigData.SerialName,
			BaudRate:     c.ConfigData.BaudRate,
			DataBits:     c.ConfigData.DataBits,
			StopBits:     c.ConfigData.StopBits,
			Parity:       c.ConfigData.Parity,
			RS485Enabled: c.ConfigData.RS485Enabled,
		}
	default:
		klog.Errorf("Invalid CommunicateMode: %s", c.ConfigData.CommunicateMode)
	}
	klog.Infoln("Start InitDevice with config:", config)
	klog.Infoln("ConfigType:", fmt.Sprintf("%T", config))
	modbusClient, err := NewModBusClient(config)
	if err != nil {
		klog.Errorf("Failed to create Modbus client: %v", err)
		return err
	}

	if err := modbusClient.Client.Connect(); err != nil {
		klog.Errorf("Failed to connect to Modbus device: %v", err)
		return err
	}
	c.ModbusClient = modbusClient
	klog.Infoln("InitDevice success")
	return nil
}

// GetDeviceData :get data from modbus client and normalize data
func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	klog.Infof("Start reading data from the device,Register: %s, Offset: %d, Limit: %d", visitor.Register, visitor.Offset, visitor.Limit)
	data, err := c.ModbusClient.Get(visitor.Register, visitor.Offset, visitor.Limit)
	if err != nil {
		klog.Errorf("Failed to read data from the device: %v", err)
		return nil, err
	}
	var value uint16
	switch visitor.Register {
	case CoilRegister, InputRegister, DiscreteInputRegister:
		//The returned Coil register data occupies 1 byte
		value = uint16(data[0])

	case HoldingRegister:
		// The returned Holding register data occupies 2 bytes
		value = binary.BigEndian.Uint16(data)
	}

	v, err := DataNormalize(value, visitor)
	return v, err
}

// ZoomIn: zoom in data
func ZoomIn(data uint16, scale float64) float64 {
	return cast.ToFloat64(data) * scale
}

// ZoomOut: zoom out data
func ZoomOut(data float64, scale float64) uint16 {
	if scale == 0 {
		return 0
	}
	return cast.ToUint16(data / scale)
}

// DataNormalize: data normalization
func DataNormalize(data interface{}, visitor *VisitorConfig) (interface{}, error) {

	switch visitor.DataType {
	case INT:
		return cast.ToInt(data), nil
	case STRING:
		return cast.ToString(data), nil
	case FLOAT, DOUBLE:
		v := cast.ToFloat64(data)

		if visitor.Scale != 0 {
			v = v * visitor.Scale
		}
		v = math.Trunc(v*1e2+0.5) * 1e-2
		return v, nil
	case BYTES:
		return data, nil
	case BOOL:
		return cast.ToBool(data), nil
	}
	return nil, fmt.Errorf("unsupported data type: %v", visitor.DataType)

}

// DeviceDataWrite: External API call to write data to the device
func (c *CustomizedClient) DeviceDataWrite(visitor *VisitorConfig, deviceMethodName string, propertyName string, data interface{}) error {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	klog.Infof("start write data to device, method name: %s, property name: %s, data: %v, register: %s", deviceMethodName, propertyName, data, visitor.Register)
	value := cast.ToUint16(ZoomOut(cast.ToFloat64(data), visitor.Scale))

	res, err := c.ModbusClient.Set(visitor.Register, visitor.Offset, value)
	if err != nil {
		klog.Errorf("fail to write to the device: %v", err)
		return err
	}
	klog.Infof("write data to device success, value: %v", binary.BigEndian.Uint16((res)))
	return nil
}

// SetDeviceData: Normalize the data for uploading
func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	v := cast.ToUint16(ZoomOut(cast.ToFloat64(data), visitor.Scale))
	visitor.Value = v
	return nil
}

// StopDevice: stop modbus client
func (c *CustomizedClient) StopDevice() error {
	err := c.ModbusClient.Client.Close()
	if err != nil {
		klog.Errorf("Failed to close Modbus client: %v", err)
		return err
	}
	return nil
}

// GetDeviceStates: check the status of device
func (c *CustomizedClient) GetDeviceStates() (string, error) {
	klog.Infoln("start to check device status")
	if err := c.ModbusClient.Client.Connect(); err != nil {
		klog.Errorf("fail to connect to device: %v", err)
		return common.DeviceStatusDisCONN, err
	}
	return common.DeviceStatusOK, nil
}
