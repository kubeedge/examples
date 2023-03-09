package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
    r.POST("/servicebus", func(c *gin.Context) {
		s, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Printf("edge say '%s' ", s)
	})
	r.Run(":8888")
}
