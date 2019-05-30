package configuration

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//CONFIG_MAP_PATH contains the location of the configuration file loaded from config map
var CONFIG_MAP_PATH = "/opt/kubeedge/deviceProfile.json"

//CONFIG_FILE_PATH contains the location of the configuration file
var CONFIG_FILE_PATH = "configuration/config.yaml"

//ReadConfigFile is the structure that is used to read the config file to get configuration information from the user
type ReadConfigFile struct {
	DeviceName string `yaml:"device-name,omitempty"`
	MQTTURL    string `yaml:"mqtt-url,omitempty"`
}

// DeviceProfile is structure to store in configMap
type DeviceProfile struct {
	DeviceInstances []DeviceInstance `json:"deviceInstances,omitempty"`
	DeviceModels    []DeviceModel    `json:"deviceModels,omitempty"`
}

// DeviceInstance is structure to store device in deviceProfile.json in configmap
type DeviceInstance struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Model    string `json:"model,omitempty"`
}

// DeviceModel is structure to store deviceModel in deviceProfile.json in configmap
type DeviceModel struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Properties  []Property `json:"properties,omitempty"`
}

// Property is structure to store deviceModel property
type Property struct {
	Name         string      `json:"name,omitempty"`
	DataType     string      `json:"dataType,omitempty"`
	Description  string      `json:"description,omitempty"`
	AccessMode   string      `json:"accessMode,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Minimum      int64       `json:"minimum,omitempty"`
	Maximum      int64       `json:"maximum,omitempty"`
	Unit         string      `json:"unit,omitempty"`
}

//ReadFromConfigMap is used to load the information from the configmaps that are provided from the cloud
func (deviceProfile *DeviceProfile) ReadFromConfigMap() error {
	jsonFile, err := ioutil.ReadFile(CONFIG_MAP_PATH)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonFile, deviceProfile)
	if err != nil {
		return err
	}
	return nil
}

//ReadFromConfigFile is used to load the information from the configuration file
func (readConfigFile *ReadConfigFile) ReadFromConfigFile() error {
	yamlFile, err := ioutil.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, readConfigFile)
	if err != nil {
		return err
	}
	return nil
}
