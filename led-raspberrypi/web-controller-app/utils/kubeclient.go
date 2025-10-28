package utils

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var KubeMaster = ""
var Kubeconfig = ""
var KubeQPS = float32(5.000000)
var KubeBurst = 10
var KubeContentType = "application/vnd.kubernetes.protobuf"

// KubeConfig from flags
func KubeConfig() (conf *rest.Config, err error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags(KubeMaster, Kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeConfig.QPS = KubeQPS
	kubeConfig.Burst = KubeBurst
	kubeConfig.ContentType = KubeContentType
	return kubeConfig, err
}
