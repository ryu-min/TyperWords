package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"words_backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *storage.Database
}

func New() *Handler {
	db := storage.New()
	return &Handler{db: db}
}

func (e *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (e *Handler) Words(c *gin.Context) {
	wordsType := c.Param("type")
	fmt.Println("request " + wordsType)
	limit := c.Query("limit")
	if len(limit) > 0 {
		fmt.Println("limit not empty")
		_, err := strconv.Atoi(limit)
		if err != nil {
			c.String(http.StatusBadRequest, "request format error, param limit should be converted to int")
			return
		} else {
			fmt.Println("limit is" + limit)
		}
	} else {
		fmt.Println("limit is empty")
	}

	wordsTypes, err := e.db.GetWordTypes()
	if err != nil {
		c.String(http.StatusBadRequest, "words service not configuring")
		return
	}

	if wordsTypes.Contains(wordsType) {
		c.String(http.StatusOK, "success request")
		return
	} else {
		c.String(http.StatusOK, "type '%s' not supported", wordsType)
		return
	}

}
