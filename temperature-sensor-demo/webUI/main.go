package main

import (
	"fmt"
	"webui/config"
	"webui/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("config:",config.DefaultConfig)
	r:=gin.Default()

	r.Static("/static","./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	api:=r.Group("/api/v1")
	{
		api.GET("/temperature",handlers.GetTemperature)
		api.GET("/switch",handlers.GetSwitch)
		api.GET("/set/temperature",handlers.SetTemperature)
		api.GET("/set/switch",handlers.SetSwitch)
	}
	r.Run(":8080")

}

