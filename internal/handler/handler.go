package handler

import (
	"fmt"
	"words_backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *storage.Database
}

func New() *Handler {
	db := storage.New()
	types, err := db.GetWordTypes()
	if err != nil {
		panic(err)
	}
	fmt.Println("types:")
	for _, t := range types {
		fmt.Println(t)
	}
	return &Handler{db: db}
}

func (e *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
