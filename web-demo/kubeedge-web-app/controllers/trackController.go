package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/kubeedge/examples/kubeedge-web-demo/kubeedge-web-app/utils"
	"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/apis/devices/v1alpha1"

	"k8s.io/client-go/rest"
)

// DeviceStatus is used to patch device status
type DeviceStatus struct {
	Status v1alpha1.DeviceStatus `json:"status"`
}

// The device id of the speaker
var deviceID = "speaker-01"

// The default namespace in which the speaker device instance resides
var namespace = "default"

// The CRD client used to patch the device instance.
var crdClient *rest.RESTClient

func init() {
	// Create a client to talk to the K8S API server to patch the device CRDs
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

// UpdateDeviceTwinWithDesiredTrack patches the desired state of
// the device twin with the track to play.
func UpdateDeviceTwinWithDesiredTrack(track string) bool {
	status := buildStatusWithDesiredTrack(track)
	deviceStatus := &DeviceStatus{Status: status}
	body, err := json.Marshal(deviceStatus)
	if err != nil {
		log.Printf("Failed to marshal device status %v", deviceStatus)
		return false
	}
	result := crdClient.Patch(utils.MergePatchType).Namespace(namespace).Resource(utils.ResourceTypeDevices).Name(deviceID).Body(body).Do()
	if result.Error() != nil {
		log.Printf("Failed to patch device status %v of device %v in namespace %v \n error:%+v", deviceStatus, deviceID, namespace, result.Error())
		return false
	} else {
		log.Printf("Track [ %s ] will be played on speaker %s", track, deviceID)
	}
	return true
}

func buildStatusWithDesiredTrack(song string) v1alpha1.DeviceStatus {
	metadata := map[string]string{"timestamp": strconv.FormatInt(time.Now().Unix()/1e6, 10),
		"type": "string",
	}
	twins := []v1alpha1.Twin{{PropertyName: "track", Desired: v1alpha1.TwinProperty{Value: song, Metadata: metadata}, Reported: v1alpha1.TwinProperty{Value: "unknown", Metadata: metadata}}}
	devicestatus := v1alpha1.DeviceStatus{Twins: twins}
	return devicestatus
}

type TrackController struct {
	BaseController
}

// Index is the initial view
func (controller *TrackController) Index() {
	log.Println("Index Start")

	controller.Layout = "layout.html"
	controller.TplName = "content.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "head.html"

	log.Println("Index Finish")
}

// PlayTrack
func (controller *TrackController) PlayTrack() {
	// Get track id
	params := struct {
		TrackID string `form:":trackId"`
	}{controller.GetString(":trackId")}

	// Validate
	if controller.ParseAndValidate(&params) == false {
		return
	}

	log.Printf("PlayTrack: %s", params.TrackID)
	// update track
	UpdateDeviceTwinWithDesiredTrack(params.TrackID)
	// response
	controller.AjaxResponse(0, "SUCCESS", nil)
}

// StopTrack
func (controller *TrackController) StopTrack() {
	log.Println("StopTrack")
	// update track
	UpdateDeviceTwinWithDesiredTrack("stop")
	// response
	controller.AjaxResponse(0, "SUCCESS", nil)
}
