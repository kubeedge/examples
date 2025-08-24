package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	IP = "127.0.0.1"
)

type DeviceData struct {
	Data    interface{} `json:"data"`
	RcdTime int         `json:"rcdTime"`
	Warning int         `json:"warning"`
}

type EdgeAPIResponse struct {
	APIVersion string `json:"apiVersion"`
	StatusCode int    `json:"statusCode"`
	TimeStamp  string `json:"timeStamp"`
	Data       struct {
		DeviceName   string      `json:"DeviceName"`
		PropertyName string      `json:"PropertyName"`
		Namespace    string      `json:"Namespace"`
		Value        string `json:"Value"` // 嵌套JSON字符串
		Type         string      `json:"Type"`
		TimeStamp    int64       `json:"TimeStamp"`
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
	if err:=json.Unmarshal([]byte(edgeResp.Data.Value),&devData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal value"+err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"temperature": devData.Data, "rcdTime": devData.RcdTime, "warning": devData.Warning})
}

func GetSwitch(c *gin.Context) {
	url := "http://" + IP + ":30077/api/v1/device/default/temperature-instance/temperature-switch"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get switch status"+err.Error()})
		return
	}
	defer resp.Body.Close()

	var edgeResp EdgeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&edgeResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body" + err.Error()})
		return
	}
	var devData DeviceData

	if err:=json.Unmarshal([]byte(edgeResp.Data.Value),&devData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal value"+err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"switch": devData.Data, "rcdTime": devData.RcdTime, "warning": devData.Warning})
}

func SetTemperature(c *gin.Context) {
	var temperature float64
	if err := c.ShouldBindJSON(&temperature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url := "http://" + IP + ":30077/api/v1/device/default/temperature-instance/temperature"
	resp,err:=http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set temperature"+err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"temperature": temperature})
}

func SetSwitch(c *gin.Context) {
	var Switch int
	if err := c.ShouldBindJSON(&Switch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url := "http://" + IP + ":30077/api/v1/device/default/temperature-instance/temperature-switch"
	resp,err:=http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set switch status"+err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"switch status": Switch})
}

