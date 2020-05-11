package main

import (
	"github.com/astaxie/beego"
	"github.com/kubeedge/examples/kubeedge-counter-demo/web-controller-app/controller"
)

func main() {
	beego.Router("/", new(controllers.TrackController), "get:Index")
	beego.Router("/track/control/:trackId", new(controllers.TrackController), "get,post:ControlTrack")

	beego.Run(":80")
}
