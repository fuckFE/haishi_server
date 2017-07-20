package server

import (
	"net/http"

	"github.com/fuckFE/haishi_server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userForm struct {
	AdminUser string `form:"adminUser" json:"adminUser"`
	AdminPass string `form:"adminPass" json:"adminPass"`
	User      string `form:"user" json:"user" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
}

func setupUsers(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("", func(c *gin.Context) {
		c.String(200, "users")
	})
	users.POST("", createUser)
	users.POST("/login", login)
}

func createUser(c *gin.Context) {
	var uf userForm
	if err := c.BindJSON(&uf); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if uf.AdminPass != "admin" || uf.AdminUser != "admin" {
		c.JSON(403, "unauth")
		return
	}

	if len(uf.User) == 0 || len(uf.Password) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "user or password can not be empty",
		})
		return
	}

	if _, err := model.GetUser(uf.User); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "user exists",
		})
		return
	}

	u, err := model.CreateUser(uf.User, uf.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

func login(c *gin.Context) {
	var uf userForm
	if err := c.BindJSON(&uf); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if len(uf.User) == 0 || len(uf.Password) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "user or password can not be empty",
		})
		return
	}

	success := model.Login(uf.User, uf.Password)

	if success {
		s := sessions.Default(c)
		s.Set("user", uf.User)
		s.Save()
		c.JSON(http.StatusOK, "ok")
		return
	}

	c.JSON(http.StatusForbidden, gin.H{
		"message": "login fail",
	})
}
