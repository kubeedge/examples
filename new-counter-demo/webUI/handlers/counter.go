package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	IP = "127.0.0.1"
)

type EdgeAPIResponse struct {
	APIVersion string `json:"apiVersion"`
	StatusCode int    `json:"statusCode"`
	TimeStamp  string `json:"timeStamp"`
	Data       struct {
		DeviceName   string      `json:"DeviceName"`
		PropertyName string      `json:"PropertyName"`
		Namespace    string      `json:"Namespace"`
		Value        string      `json:"Value"` 
		Type         string      `json:"Type"`
		TimeStamp    int64       `json:"TimeStamp"`
	} `json:"Data"`
}

func GetCounter(c *gin.Context) {
	url := "http://" + IP + ":30077/api/v1/device/default/counter-instance/count"
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
	fmt.Println(edgeResp)

	c.JSON(http.StatusOK, gin.H{"counter": edgeResp.Data.Value})
}

func GetStatus(c *gin.Context) {
	url := "http://" + IP + ":30077/api/v1/device/default/counter-instance/status"
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

	c.JSON(http.StatusOK, gin.H{"status": edgeResp.Data.Value})
}


func SetStatus(c *gin.Context) {
	Status :=c.Query("status")

	url := "http://" + IP + ":30077/api/v1/devicemethod/default/counter-instance/UpdateStatus/status/"+Status
	resp,err:=http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set switch status"+err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"status": Status, "success": true})
}

func ResetCount(c *gin.Context) {
	url := "http://" + IP + ":30077/api/v1/devicemethod/default/counter-instance/SetCount/count/0"
	resp,err:=http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset count"+err.Error()})
		return
	}
	defer resp.Body.Close()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func SetCount(c *gin.Context) {
	Count :=c.Query("count")
	url := "http://" + IP + ":30077/api/v1/devicemethod/default/counter-instance/SetCount/count/"+Count
	resp,err:=http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set count"+err.Error()})
		return
	}
	defer resp.Body.Close()
	c.JSON(http.StatusOK, gin.H{"count": Count, "success": true})
}
