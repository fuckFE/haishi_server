package server

import (
	"github.com/gin-gonic/gin"
)

func SetupAPI(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.String(200, "haishi server")
	})
}
