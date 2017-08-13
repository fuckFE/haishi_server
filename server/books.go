package server

import (
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/fuckFE/haishi_server/model"
	"github.com/gin-gonic/gin"
)

type bookForm struct {
	Types    []int  `form:"types" json:"types" binding:"required"`
	Filename string `form:"filename" json:"filename" binding:"required"`
	FileID   int    `form:"fileId" json:"fileId" binding:"required"`
}

func setupBooks(rg *gin.RouterGroup) {
	books := rg.Group("/books")
	books.POST("", createBook)
	books.GET("", getBookByTypeId)
	books.POST("/grabate", getGarbate)
	books.GET("/:bookID", getBookByID)
	books.PUT("/:bookID/garbate", garbate)
	books.PUT("/:bookID/payload", payload)
	books.DELETE("/:bookID", delBook)
}

func createBook(c *gin.Context) {
	var bf bookForm
	if err := c.BindJSON(&bf); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tf, err := model.GetTmpfileById(uint64(bf.FileID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	types := make([]uint64, 0)
	for _, t := range bf.Types {
		if t > 0 {
			types = append(types, uint64(t))
		}
	}

	tf.Filename = strings.TrimSuffix(tf.Filename, path.Ext(tf.Filename))
	if err := model.CreateBook(tf.Filename, tf.Payload, types); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func getGarbate(c *gin.Context) {
	bs, err := model.GetBookByGrabate()
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, bs)
}

func getBookByTypeId(c *gin.Context) {
	id := c.Query("type")
	filterGarbateStr := c.Query("filterGarbate")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}

	filterGarbate := false
	if len(filterGarbateStr) > 0 {
		filterGarbate = true
	}
	bs, err := model.GetBook(uint64(uid), filterGarbate)
	for _, b := range bs {
		b.Payload = ""
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}

	c.JSON(200, bs)
}

func garbate(c *gin.Context) {
	id := c.Param("bookID")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}

	if err := model.Garbate(uint64(uid)); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(200, "ok")
}

func payload(c *gin.Context) {
	id := c.Param("bookID")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fileID := c.PostForm("fileID")
	uFileID, err := strconv.Atoi(fileID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tf, err := model.GetTmpfileById(uint64(uFileID))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := model.UpdatePayload(uint64(uid), tf.Payload, tf.Filename); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	model.DelTmpFileById(tf.ID)
	c.String(200, "OK")
}

func delBook(c *gin.Context) {
	id := c.Param("bookID")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}

	if err := model.DelBook(uint64(uid)); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(200, "ok")
}

func getBookByID(c *gin.Context) {
	id := c.Param("bookID")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}

	b, err := model.GetBookByID(uint64(uid))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, b)
}
