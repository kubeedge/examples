package controllers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/astaxie/beego"
	"github.com/kubeedge/examples/kubeedge-counter-demo/web-controller-app/utils"
	"k8s.io/client-go/rest"
)

var deviceID = "counter"
var namespace = "default"
var originCmd = "OFF"
var crdClient *rest.RESTClient

func init() {
	kubeConfig, err := utils.KubeConfig()
	if err != nil {
		log.Fatalf("Failed to create KubeConfig, error : %v", err)
	}
	log.Println("Get kubeConfig successfully")

	crdClient, err = utils.NewCRDClient(kubeConfig)
	if err != nil {
		log.Fatalf("Failed to create device crd client , error : %v", err)
	}
	log.Println("Get crdClient successfully")
}

// UpdateStatus reads desired + reported state from v1beta1 Device CR.
func UpdateStatus() map[string]string {
	result := struct {
		Status struct {
			Twins []struct {
				PropertyName string `json:"propertyName"`
				Reported     struct { Value string `json:"value"`} `json:"reported"`
				ObservedDesired struct { Value string `json:"value"`} `json:"observedDesired"`
			} `json:"twins"`
		} `json:"status"`
	}{}

	raw, err := crdClient.Get().
		Namespace(namespace).
		Resource(utils.ResourceTypeDevices).
		Name(deviceID).
		DoRaw(context.TODO())
	if err != nil {
		log.Printf("Failed to get device %v: %v", deviceID, err)
	}

	status := map[string]string{
		"status": "OFF",
		"value":  "0",
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		log.Printf("Failed to unmarshal device status: %v", err)
		return status
	}

	for _, twin := range result.Status.Twins {
		if twin.PropertyName == "status" {
			if twin.ObservedDesired.Value != "" {
				status["status"] = twin.ObservedDesired.Value
			}
			if twin.Reported.Value != "" {
				status["value"] = twin.Reported.Value
			}
		}
	}

	return status
}

// DeviceSpecPatch patches spec.properties desired value for v1beta1.
type DeviceSpecPatch struct {
	Spec struct {
		Properties []DeviceProperty `json:"properties"`
	} `json:"spec"`
}

type DeviceProperty struct {
	Name    string `json:"name"`
	Desired struct {
		Value string `json:"value"`
	} `json:"desired"`
}

func buildDeviceSpecPatch(cmd string) *DeviceSpecPatch {
	patch := &DeviceSpecPatch{}

	prop := DeviceProperty{
		Name: "status",
	}
	prop.Desired.Value = cmd

	patch.Spec.Properties = []DeviceProperty{prop}
	return patch
}

func UpdateDeviceTwinWithDesiredTrack(cmd string) bool {
	if cmd == originCmd {
		return true
	}

	devicePatch := buildDeviceSpecPatch(cmd)

	body, err := json.Marshal(devicePatch)
	if err != nil {
		log.Printf("Failed to marshal device patch %v", devicePatch)
		return false
	}

	result := crdClient.Patch(utils.MergePatchType).
		Namespace(namespace).
		Resource(utils.ResourceTypeDevices).
		Name(deviceID).
		Body(body).
		Do(context.TODO())

	if result.Error() != nil {
		log.Printf("Failed to patch device %v in namespace %v error:%+v",
			deviceID, namespace, result.Error())
		return false
	}

	log.Printf("Turn %s %s", cmd, deviceID)
	originCmd = cmd
	return true
}

type TrackController struct {
	beego.Controller
}

func (controller *TrackController) Index() {
	log.Println("Index Start")
	controller.Layout = "layout.html"
	controller.TplName = "content.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "head.html"
	log.Println("Index Finish")
}

func (controller *TrackController) ControlTrack() {
	params := struct {
		TrackID string `form:":trackId"`
	}{controller.GetString(":trackId")}

	resultCode := 0
	status := map[string]string{}
	log.Printf("ControlTrack: %s", params.TrackID)

	if params.TrackID == "ON" {
		UpdateDeviceTwinWithDesiredTrack(params.TrackID)
		resultCode = 1
	} else if params.TrackID == "OFF" {
		UpdateDeviceTwinWithDesiredTrack(params.TrackID)
		resultCode = 2
	} else if params.TrackID == "STATUS" {
		status = UpdateStatus()
		resultCode = 3
	}

	controller.AjaxResponse(resultCode, status, nil)
}

func (Controller *TrackController) AjaxResponse(resultCode int, resultString map[string]string, data interface{}) {
	response := struct {
		Result       int
		ResultString map[string]string
		ResultObject interface{}
	}{
		Result:       resultCode,
		ResultString: resultString,
		ResultObject: data,
	}

	Controller.Data["json"] = response
	Controller.ServeJSON()
}
