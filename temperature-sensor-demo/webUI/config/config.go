package config

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cast"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
type DeviceCustomConfig struct {
	MaxT float64 `json:"maxT"`
	MinT float64 `json:"minT"`
}

var DefaultConfig DeviceCustomConfig 

const(
	DeviceName = "temperature-instance"
)

func init(){
	GetDefaultConfig()
}

func GetDefaultConfig(){
	item := GetCRD()
	if item == nil{
		fmt.Println("Failed to find target device")
		return
	}
	fmt.Println(item.GetName(),"spec:\n",item.Object["spec"])
	properties, _, err := unstructured.NestedFieldNoCopy(item.Object, "spec", "properties")
	if err!= nil{
		fmt.Printf("Failed to get properties")
	}
	var maxT, minT interface{}
	for _, property := range properties.([]interface{}) {
		name, _, _ := unstructured.NestedFieldNoCopy(property.(map[string]interface{}), "name")
		if name == "temperature"{
			maxT, _, _ = unstructured.NestedFieldNoCopy(property.(map[string]interface{}), "visitors","configData","max")
			minT, _, _ = unstructured.NestedFieldNoCopy(property.(map[string]interface{}), "visitors","configData","min")
		}
	}
	DefaultConfig.MaxT = cast.ToFloat64(maxT)
	DefaultConfig.MinT = cast.ToFloat64(minT)
}
// auto get kubeconfig
func getRestConfig() (*rest.Config, error) {
	
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = home + "/.kube/config"
	}
	
	if _, err := os.Stat(kubeconfig); err == nil {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func GetCRD() *unstructured.Unstructured{
	config, err := getRestConfig()
	if err != nil || config == nil {
		fmt.Println("Failed to get kubeconfig")
		return nil
	}
	fmt.Println("kubeconfig:",config)

	// create dynamic client
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Println("Failed to create dynamic client")
		return nil
	}
	// Define CRD GVR
	gvr := schema.GroupVersionResource{
		Group:    "devices.kubeedge.io",
		Version:  "v1beta1",
		Resource: "devices",
	}

	// find device instance in default namespace
	list, err := dyn.Resource(gvr).Namespace("default").List(context.Background(), 
		metav1.ListOptions{})
	if err != nil {
		fmt.Println("Failed to get CRD")
		return nil
	}

	fmt.Println("Device CRD Instance List:")
	for _, item := range list.Items {
		fmt.Println(item.GetName())
		if item.GetName()==DeviceName{
			return &item
		}	
	}
	return nil
}