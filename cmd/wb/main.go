package main

import (
	_ "embed"
	"words_backend/internal/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	h := handler.New()
	r := gin.Default()
	r.GET("/ping", h.Ping)
	r.Run() // listen and serve on 0.0.0.0:8080
}
