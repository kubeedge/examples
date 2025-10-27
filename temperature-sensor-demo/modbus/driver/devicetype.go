package driver

import (
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
)

const (
	CoilRegister          = "CoilRegister"
	InputRegister         = "InputRegister"
	HoldingRegister       = "HoldingRegister"
	DiscreteInputRegister = "DiscreteInputRegister"
	INT                   = "int"
	FLOAT                 = "float"
	DOUBLE                = "double"
	STRING                = "string"
	BOOL                  = "bool"
	BYTES                 = "bytes"
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
	ModbusClient *ModbusClient
}

type ProtocolConfig struct {
	ProtocolName string `json:"protocolName"`
	ConfigData   `json:"configData"`
}

type ConfigData struct {
	SlaveID         int    `json:"slaveID,omitempty"`
	CommunicateMode string `json:"communicateMode"` // enum:TCP/RTU
	// TCP mode
	IP      string `json:"ip,omitempty"`
	Port    string `json:"port,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
	// RTU mode
	SerialName   string `json:"serialName,omitempty"`
	BaudRate     int    `json:"baudRate,omitempty"`
	DataBits     int    `json:"dataBits,omitempty"`
	StopBits     int    `json:"stopBits,omitempty"`
	Parity       string `json:"parity,omitempty"` // enum:None/Even/Odd
	RS485Enabled bool   `json:"rs485Enabled,omitempty"`
}

type VisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	DataType       string      `json:"dataType"` // enum:Int/Float/String
	Value          interface{} `json:"value"`
	Register       string      `json:"register"` // enum:CoilRegister/HoldingRegister
	Offset         uint16      `json:"offset"`
	Scale          float64     `json:"scale"`
	IsSwap         bool        `json:"isSwap,omitempty"`
	Limit          uint16      `json:"limit"`
	IsRegisterSwap bool        `json:"isRegisterSwap,omitempty"`
}