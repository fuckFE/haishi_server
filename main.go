package main

import (
	"github.com/fuckFE/haishi_server/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})
	server.SetupAPI(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
