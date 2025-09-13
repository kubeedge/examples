package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webui/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const (
	IP = "127.0.0.1"
)

type DeviceData struct {
	Data    interface{} `json:"data"`
	Warning int         `json:"warning"`
}

type EdgeAPIResponse struct {
	APIVersion string `json:"apiVersion"`
	StatusCode int    `json:"statusCode"`
	TimeStamp  string `json:"timeStamp"`
	Data       struct {
		DeviceName   string `json:"DeviceName"`
		PropertyName string `json:"PropertyName"`
		Namespace    string `json:"Namespace"`
		Value        string `json:"Value"`
		Type         string `json:"Type"`
		TimeStamp    int64  `json:"TimeStamp"`
	} `json:"Data"`
}

func GetTemperature(c *gin.Context) {
	url := "http://" + IP + ":30077/api/v1/device/default/temperature-instance/temperature"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var edgeResp EdgeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&edgeResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body" + err.Error()})
		return
	}
	var devData DeviceData

	devData.Data = edgeResp.Data.Value
	devData.Warning = getWarning(cast.ToFloat64(devData.Data))

	c.JSON(http.StatusOK, gin.H{"temperature": devData.Data, "warning": devData.Warning})
}

func GetSwitch(c *gin.Context) {
	url := "http://" + IP + ":30077/api/v1/device/default/temperature-instance/temperature-switch"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get switch status" + err.Error()})
		return
	}
	defer resp.Body.Close()

	var edgeResp EdgeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&edgeResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body" + err.Error()})
		return
	}
	var devData DeviceData
	devData.Data = edgeResp.Data.Value
	c.JSON(http.StatusOK, gin.H{"switch": devData.Data})
}

func getWarning(t float64) int {
	fmt.Println(config.DefaultConfig)
	if t > config.DefaultConfig.MaxT {
		return 1
	}
	if t < config.DefaultConfig.MinT {
		return 2
	}
	return 0
}

func SetTemperature(c *gin.Context) {
	temperature := c.Query("temperature")
	url := "http://" + IP + ":30077/api/v1/devicemethod/default/temperature-instance/UpdateTemperature/temperature/" + temperature
	resp, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set temperature" + err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"temperature": temperature, "success": true})
}

func SetSwitch(c *gin.Context) {
	SwitchStatus := c.Query("switchStatus")
	url := "http://" + IP + ":30077/api/v1/devicemethod/default/temperature-instance/SwitchControl/temperature-switch/" + SwitchStatus
	resp, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set switch status" + err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"switch status": SwitchStatus, "success": true})
}
