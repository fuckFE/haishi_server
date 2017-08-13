package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fuckFE/haishi_server/model"
	"github.com/gin-gonic/gin"
)

type typeForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Category uint   `form:"category" json:"category" binding:"required"`
}

func setupTypes(rg *gin.RouterGroup) {
	types := rg.Group("/types")
	types.GET("", getTypes)
	types.POST("", createType)
	types.DELETE("/:id", delTypeById)
}

func getTypes(c *gin.Context) {
	ts, err := model.GetTypes()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, ts)
}

func createType(c *gin.Context) {
	var tf typeForm
	if err := c.BindJSON(&tf); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	t, err := model.CreateType(tf.Name, tf.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, t)
}

func delTypeById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	if err := model.DelTypeById(uint64(id)); err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, "success")
}
