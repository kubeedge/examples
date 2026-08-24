package http

import (
	"testing"

	"github.com/kubeedge/mapper-framework/pkg/common"
)

func TestPushReturnsWhenHTTPPostFails(t *testing.T) {
	pushMethod := &PushMethod{HTTP: &HTTPConfig{
		HostName:    "://invalid",
		RequestPath: "/publish",
	}}

	pushMethod.Push(&common.DataModel{PropertyName: "count", Value: "1"})
}
