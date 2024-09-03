package control

import (
	"encoding/json"
	"errors"
	"net/http"

	"k8s.io/klog/v2"

	"github.com/kubeedge/examples/sound-equipment-fault-detection/vled/device"
)

type RequestData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func Send() {
	http.HandleFunc("/api/v1/resource", handler)
	http.ListenAndServe(":5050", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var reqData RequestData

	// Parsing JSON request body
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process Request Data
	rdata, err := ProcessRequestData(reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Set Observed Desired Value
	device.SetObservedDesiredValue(rdata)

	// Print received data
	klog.Infoln("The controller has received the data")
	klog.Infoln(reqData.Type, reqData.Data)

	// Return Response
	w.Write([]byte("Data received successfully!"))
}

func ProcessRequestData(reqData RequestData) (string, error) {
	rtype := reqData.Type
	rdata := reqData.Data
	if rtype == "set vled" {
		return rdata, nil
	} else {
		return "", errors.New("rtype cannot be resolved")
	}
}
