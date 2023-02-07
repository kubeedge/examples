module github.com/kubeedge/examples/kubeedge-counter-demo

go 1.14

replace (
	github.com/apache/servicecomb-kie v0.1.0 => github.com/apache/servicecomb-kie v0.0.0-20190905062319-5ee098c8886f // indirect. TODO: remove this line when servicecomb-kie has a stable release
	github.com/kubeedge/beehive => github.com/kubeedge/beehive v0.0.0-20201125122335-cd19bca6e436 // indirect
	github.com/kubeedge/viaduct => github.com/kubeedge/viaduct v0.0.0-20201130063818-e33931917980 // indirect
	k8s.io/api v0.0.0 => k8s.io/api v0.0.0-20190720062849-3043179095b6
	k8s.io/apiextensions-apiserver v0.0.0 => k8s.io/apiextensions-apiserver v0.0.0-20190718185103-d1ef975d28ce // indirect
	k8s.io/apimachinery v0.0.0 => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/apiserver v0.0.0 => k8s.io/apiserver v0.0.0-20190718184206-a1aa83af71a7
	k8s.io/cli-runtime v0.0.0 => k8s.io/cli-runtime v0.0.0-20190718185405-0ce9869d0015
	k8s.io/client-go v0.0.0 => k8s.io/client-go v0.0.0-20190718183610-8e956561bbf5 // indirect
	k8s.io/cloud-provider v0.0.0 => k8s.io/cloud-provider v0.0.0-20190718190308-f8e43aa19282 // indirect
	k8s.io/cluster-bootstrap v0.0.0 => k8s.io/cluster-bootstrap v0.0.0-20190718190146-f7b0473036f9
	k8s.io/code-generator v0.0.0 => k8s.io/code-generator v0.15.8-beta.1
	k8s.io/component-base v0.0.0 => k8s.io/component-base v0.0.0-20190718183727-0ececfbe9772
	k8s.io/cri-api v0.0.0 => k8s.io/cri-api v0.0.0-20190531030430-6117653b35f1
	k8s.io/csi-api v0.0.0 => k8s.io/csi-api v0.0.0-20190313123203-94ac839bf26c // indirect
	k8s.io/csi-translation-lib v0.0.0 => k8s.io/csi-translation-lib v0.0.0-20190718190424-bef8d46b95de
	k8s.io/gengo v0.0.0 => k8s.io/gengo v0.0.0-20190327210449-e17681d19d3a // indirect
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1 // indirect
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.2.0
	k8s.io/kube-aggregator v0.0.0 => k8s.io/kube-aggregator v0.0.0-20190718184434-a064d4d1ed7a
	k8s.io/kube-controller-manager v0.0.0 => k8s.io/kube-controller-manager v0.0.0-20190718190030-ea930fedc880
	k8s.io/kube-openapi v0.0.0 => k8s.io/kube-openapi v0.0.0-20190718094010-3cf2ea392886 // indirect
	k8s.io/kube-proxy v0.0.0 => k8s.io/kube-proxy v0.0.0-20190718185641-5233cb7cb41e
	k8s.io/kube-scheduler v0.0.0 => k8s.io/kube-scheduler v0.0.0-20190718185913-d5429d807831
	k8s.io/kubectl => k8s.io/kubectl v0.19.1
	k8s.io/kubelet v0.0.0 => k8s.io/kubelet v0.0.0-20190718185757-9b45f80d5747
	k8s.io/legacy-cloud-providers v0.0.0 => k8s.io/legacy-cloud-providers v0.0.0-20190718190548-039b99e58dbd
	k8s.io/metrics v0.0.0 => k8s.io/metrics v0.0.0-20190718185242-1e1642704fe6
	k8s.io/node-api v0.0.0 => k8s.io/node-api v0.0.0-20190717025432-9e6fdeee55cc // indirect
	k8s.io/repo-infra v0.0.0 => k8s.io/repo-infra v0.0.0-20181204233714-00fe14e3d1a3 // indirect
	k8s.io/sample-apiserver v0.0.0 => k8s.io/sample-apiserver v0.0.0-20190718184639-baafa86838c0
	k8s.io/utils v0.0.0 => k8s.io/utils v0.0.0-20190712204705-3dccf664f023
)

require (
	github.com/astaxie/beego v1.12.3
	github.com/eclipse/paho.mqtt.golang v1.3.1
	github.com/kubeedge/kubeedge v1.5.0
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.20.0
)
