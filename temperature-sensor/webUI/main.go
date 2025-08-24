package main

import (
	"webui/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r:=gin.Default()

	r.Static("/static","./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	api:=r.Group("/api/v1")
	{
		api.GET("/temperature",handlers.GetTemperature)
		api.GET("/switch",handlers.GetSwitch)
		api.POST("/set/temperature",handlers.SetTemperature)
		api.POST("/set/switch",handlers.SetSwitch)
	}
	r.Run(":8080")
}

