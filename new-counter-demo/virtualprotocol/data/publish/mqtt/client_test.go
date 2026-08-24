package mqtt

import (
	"testing"

	"github.com/kubeedge/mapper-framework/pkg/common"
)

func TestPushReturnsWhenMQTTConnectionFails(t *testing.T) {
	pushMethod := &PushMethod{MQTT: &MQTTConfig{Address: "://"}}
	pushMethod.Push(&common.DataModel{Value: "1"})
}
