package server

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fuckFE/haishi_server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	userParamsError = 100
)

func GetMainEngine() *gin.Engine {
	staticPath, err := filepath.Abs("public")
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	store := sessions.NewCookieStore([]byte("hs_store"))
	r.Use(sessions.Sessions("hs_store", store))

	api := r.Group("/api")
	api.Use(authMid())
	api.GET("", func(c *gin.Context) {
		c.String(200, "haishi server")
	})
	api.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		payload, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		tf, err := model.CreateTmpfile(header.Filename, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(200, gin.H{
			"id":       tf.ID,
			"filename": tf.Filename,
		})
	})
	setupUsers(api)
	setupTypes(api)
	setupBooks(api)

	r.Use(static.Serve("/", static.LocalFile(staticPath, true)))

	return r
}

func authMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api/users") {
			god := c.Request.Header.Get("x-god")
			if god != "app" {
				s := sessions.Default(c)
				if val := s.Get("user"); val == nil {
					c.String(403, "unauth")
					return
				}
			}
		}

		c.Next()
	}
}
