package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kubeedge/mapper-framework/pkg/common"
	externalanomalydetectorpb "github.com/kubeedge/me-restful-mapper/api/externalanomalydetectorpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/klog/v2"
)

func NewClient(protocol ProtocolConfig) (*CustomizedClient, error) {
	client := &CustomizedClient{
		ProtocolConfig: protocol,
		deviceMutex:    sync.Mutex{},
	}
	return client, nil
}

func (c *CustomizedClient) InitDevice() error {
	return nil
}

func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	ret, err := GetDeviceValue(c.ProtocolConfig.ConfigData.Ip, c.ProtocolConfig.ConfigData.Port, c.ProtocolConfig.ConfigData.Url)
	if err != nil {
		klog.Errorf("GetDeviceData failed for http://%s:%d%s, error: %v", c.ProtocolConfig.ConfigData.Ip, c.ProtocolConfig.ConfigData.Port, c.ProtocolConfig.ConfigData.Url, err)
		return nil, nil
	}
	klog.V(4).Infof("GetDeviceData success for http://%s:%d%s, data: %s", c.ProtocolConfig.ConfigData.Ip, c.ProtocolConfig.ConfigData.Port, c.ProtocolConfig.ConfigData.Url, ret)
	return ret, nil
}

func (c *CustomizedClient) DeviceDataWrite(visitor *VisitorConfig, deviceMethodName string, propertyName string, data interface{}) error {
	return nil
}

func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	err := SetDeviceValue(c.ProtocolConfig.ConfigData.Ip, c.ProtocolConfig.ConfigData.Port, c.ProtocolConfig.ConfigData.Url, data.(string))
	if err != nil {
		klog.Errorf("DeviceDataWrite failed for http://%s:%d%s, error: %v", c.ProtocolConfig.ConfigData.Ip, c.ProtocolConfig.ConfigData.Port, c.ProtocolConfig.ConfigData.Url, err)
		return nil
	}
	klog.V(4).Infof("DeviceDataWrite success for http://%s:%d%s, data: %s", c.ProtocolConfig.ConfigData.Ip, c.ProtocolConfig.ConfigData.Port, c.ProtocolConfig.ConfigData.Url, data.(string))
	return nil
}

func (c *CustomizedClient) StopDevice() error {
	return nil
}

func (c *CustomizedClient) GetDeviceStates() (string, error) {
	return common.DeviceStatusOK, nil
}

var connectionManager = NewGRPCConnectionManager()

func (c *CustomizedClient) AnomalyDetectionProcess(req *AnomalyDetectionRequest) error {
	var config AnomalyDetectionConfig
	err := json.Unmarshal(req.AnomalyDetectionConfig, &config)
	if err != nil {
		klog.Errorf("AnomalyDetectionProcess unmarshal anomalyDetectionConfig failed, error: %v", err)
		return err
	}

	// Use connection manager to get or create a gRPC connection
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := connectionManager.GetConnection(address)
	if err != nil {
		klog.Errorf("Failed to get gRPC connection to %s, error: %v", address, err)
		return err
	}

	client := externalanomalydetectorpb.NewExternalAnomalyDetectorClient(conn)

	grpcRequest := &externalanomalydetectorpb.AnomalyRequest{
		Timestamp:    timestamppb.New(time.Now()),
		DeviceName:   config.DeviceName,
		Namespace:    config.Namespace,
		PropertyName: config.PropertyName,
		PropertyValue: req.Data.(string),
	}

	// Call the gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	response, err := client.DetectAnomaly(ctx, grpcRequest)
	if err != nil {
		klog.Errorf("Failed to call DetectAnomaly, error: %v", err)
		return err
	}

	if !response.Success {
		klog.Errorf("Anomaly detection failed, errorCode: %d, message: %s", response.ErrorCode, response.Message)
		return fmt.Errorf("anomaly detection failed: %s", response.Message)
	}

	klog.V(4).Infof("Anomaly detection succeeded, message: %s", response.Message)
	return nil
}
