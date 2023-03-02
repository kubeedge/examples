package main

import (
	"github.com/astaxie/beego"
	controllers "github.com/kubeedge/examples/led-raspberrypi/web-controller-app/controller"
)

func main() {
	beego.Router("/", new(controllers.TrackController), "get:Index")
	beego.Router("/track/control/:trackId", new(controllers.TrackController), "get,post:ControlTrack")

	beego.Run(":80")
}
