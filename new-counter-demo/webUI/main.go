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
		api.GET("/counter",handlers.GetCounter)
		api.GET("/status",handlers.GetStatus)
		api.GET("/set/status",handlers.SetStatus)
		api.GET("/set/counter",handlers.SetCount)
		api.GET("/reset",handlers.ResetCount)
	}
	r.Run(":8080")
}

