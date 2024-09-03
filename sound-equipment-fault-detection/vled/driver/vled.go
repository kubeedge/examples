package driver

import (
	"sync"

	"k8s.io/klog/v2"
)

const (
	LEDStatusOff int = 0
	LEDStatusOn  int = 1
)

var (
	vLED = LEDStatusOff // The default state is closed
	mu   sync.Mutex
)

// Setting LED Status
func SetVLEDStatus(status int) {
	mu.Lock()
	defer mu.Unlock()
	if status == LEDStatusOn || status == LEDStatusOff {
		vLED = status
	} else {
		klog.Infoln("Invalid status")
	}
}

// Get LED status
func GetVLEDStatus() int {
	mu.Lock()
	defer mu.Unlock()
	return vLED
}
