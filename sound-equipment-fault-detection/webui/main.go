package main

import (
	"github.com/gin-gonic/gin"

	"github.com/kubeedge/examples/sound-equipment-fault-detection/webui/controllers"
)

func main() {

	r := gin.Default()

	// Serve static files from the 'static' directory
	r.Static("/static", "./static")

	// Serve index.html on the root route
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Use the handler function for the route
	r.GET("/api/led", controllers.GetLEDStatus)

	// Start the server
	r.Run(":8080")
}
