package utils

import (
	"log"

	devices "github.com/kubeedge/kubeedge/cloud/pkg/apis/devices/v1alpha2"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

// NewCRDClient is used to create a restClient for crd
func NewCRDClient(cfg *rest.Config) (*rest.RESTClient, error) {
	scheme := runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(addDeviceCrds)

	err := schemeBuilder.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	config := *cfg
	config.APIPath = "/apis"
	config.GroupVersion = &devices.SchemeGroupVersion
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		log.Fatalf("Failed to create REST Client due to error %v", err)
		return nil, err
	}

	return client, nil
}

func addDeviceCrds(scheme *runtime.Scheme) error {
	// Add Device
	scheme.AddKnownTypes(devices.SchemeGroupVersion, &devices.Device{}, &devices.DeviceList{})
	v1.AddToGroupVersion(scheme, devices.SchemeGroupVersion)
	// Add DeviceModel
	scheme.AddKnownTypes(devices.SchemeGroupVersion, &devices.DeviceModel{}, &devices.DeviceModelList{})
	v1.AddToGroupVersion(scheme, devices.SchemeGroupVersion)

	return nil
}
